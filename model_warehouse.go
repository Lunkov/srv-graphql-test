package main

import (
  "sync"
  "encoding/gob"
  "os"
  "unsafe"
  "errors"
  "github.com/golang/glog"
  "github.com/graphql-go/graphql"
)

// Склады

type Warehouse struct {
  CODE         string    `db:"code"           json:"code"            yaml:"code"`
  Region_ID    int64  
  Name         string    `db:"name"           json:"name"            yaml:"name"`
  Description  string    `db:"description"    json:"description"     yaml:"description"`
}

var warehouseType = graphql.NewObject(
   graphql.ObjectConfig{
      Name: "Warehouse",
      Fields: graphql.Fields{
         "code": &graphql.Field{
            Type: graphql.String,
         },
         "region_id": &graphql.Field{
            Type: graphql.String,
         },
         "name": &graphql.Field{
            Type: graphql.String,
         },
         "description": &graphql.Field{
            Type: graphql.String,
         },
      },
   },
)

var memWH map[string]Warehouse
var memWHReg = make(map[int64][]string, 100) // index region_id -> warehouses
var muWH   sync.RWMutex

func WarehouseInit(max int64) {
  memWH = make(map[string]Warehouse, max)
  glog.Infof("LOG: Warehouses: max          = %d", max)
  glog.Infof("LOG: Warehouses: sizeof(item) = %d", unsafe.Sizeof(Warehouse{}))
  glog.Infof("LOG: Warehouses: sizeof(map)  = %d", unsafe.Sizeof(memWH))
}

func WarehouseCount() int64 {
  return int64(len(memWH))
}


func GetWarehouses(offset int, limit int) []Warehouse {
  res := make([]Warehouse, limit)
  i := 0
  for _, item := range memWH {
    if offset <= i &&  offset + limit > i {
      res = append(res, item)
    }
    if offset + limit < i {
      break
    }
    i++
  }
  glog.Infof("LOG: Shops %v", res)
  return res
}

func WarehouseAppend(info *Warehouse) {
  muWH.Lock()
  memWH[info.CODE] = *info
  if _, ok := memWHReg[info.Region_ID]; !ok {
    memWHReg[info.Region_ID] = make([]string, 1)
  }
  memWHReg[info.Region_ID] = append(memWHReg[info.Region_ID], info.CODE)
  muWH.Unlock()
}

func GetWarehouseByCode(code string) (*Warehouse) {
  item, ok := memWH[code]
  if ok {
    return &item
  }
  return nil
}

func GetWarehousesByRegionID(region_id int64) ([]string) {
  muWH.RLock()
  items, ok := memWHReg[region_id]
  muWH.RUnlock()
  if ok {
    return items
  }
  return make([]string, 0)
}

func WriteFileWarehouses(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, _ := os.Create(filename)
  defer file.Close()
  encoder := gob.NewEncoder(file)
  encoder.Encode(memWH)

  fileIndex, _ := os.Create(filename+".index")
  defer fileIndex.Close()
  encoderIndex := gob.NewEncoder(fileIndex)
  encoderIndex.Encode(memWHReg)
}


func LoadFileWarehouses(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, err := os.Open(filename)
  if err !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename, err)
    return
  }
  defer file.Close()
  
  decoder := gob.NewDecoder(file)
  err = decoder.Decode(&memWH)
  if err != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename, err)
    return
  }
  
  fileI, errI := os.Open(filename+".index")
  if errI !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename+".index", errI)
    return
  }
  defer file.Close()
  
  decoderI := gob.NewDecoder(fileI)
  errI = decoderI.Decode(&memWHReg)
  if errI != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename+".index", errI)
    return
  }
}


func WarehouseGQL() {
  AppendFields2GraphQL("warehouse", &graphql.Field{
			Type: warehouseType,
      Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                  Description: "id of the warehouses",
                  Type:graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id, ok := p.Args["id"].(string)
        if !ok {
          return nil, errors.New("Need ID")
        }
				return GetWarehouseByCode(id), nil
			},
		})
    
	AppendFields2GraphQL("warehouses", &graphql.Field{
			Type: graphql.NewList(warehouseType),
      Args: graphql.FieldConfigArgument{
               "offset": &graphql.ArgumentConfig{
                  Type: graphql.Int,
               },
               "limit": &graphql.ArgumentConfig{
                  Type: graphql.Int,
               },
            },      
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        offset, ok := p.Args["offset"].(int)
        if !ok {
          offset = 0
        }
        limit, ok := p.Args["limit"].(int)
        if !ok {
          limit = 1000
        }
        glog.Infof("LOG: Warehouses(offset=%d, limit=%d)", offset, limit)
				return GetWarehouses(offset, limit), nil
			},
    })
}
