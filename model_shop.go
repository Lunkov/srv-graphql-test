package main

import (
  "sync"
  "encoding/gob"
  "os"
  "errors"
  "unsafe"
  "github.com/golang/glog"
  "github.com/graphql-go/graphql"
  
)

// Магазины

type Shop struct {
  CODE         string    `db:"code"           json:"code"            yaml:"code"`
  Region_ID    int64  
  Name         string    `db:"name"           json:"name"            yaml:"name"`
  Description  string    `db:"description"    json:"description"     yaml:"description"`
}

var shopType = graphql.NewObject(
   graphql.ObjectConfig{
      Name: "Shop",
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

var memShop map[string]Shop
var memShopReg = make(map[int64][]string, 100) // index region_id -> shops
var muShop   sync.RWMutex

func ShopInit(max int64) {
  memShop = make(map[string]Shop, max)
  glog.Infof("LOG: Shops: max          = %d", max)
  glog.Infof("LOG: Shops: sizeof(item) = %d", unsafe.Sizeof(Shop{}))
  glog.Infof("LOG: Shops: sizeof(map)  = %d", unsafe.Sizeof(memShop))
}

func ShopCount() int {
  return len(memShop)
}

func ShopAppend(info *Shop) {
  muShop.Lock()
  memShop[info.CODE] = *info
  if _, ok := memShopReg[info.Region_ID]; !ok {
    memShopReg[info.Region_ID] = make([]string, 1)
  }
  memShopReg[info.Region_ID] = append(memShopReg[info.Region_ID], info.CODE)
  muShop.Unlock()
}

func GetShopByID(code string) (*Shop) {
  muShop.RLock()
  item, ok := memShop[code]
  muShop.RUnlock()
  if ok {
    return &item
  }
  return nil
}

func GetShopByRegionID(region_id int64) ([]string) {
  muShop.RLock()
  items, ok := memShopReg[region_id]
  muShop.RUnlock()
  if ok {
    return items
  }
  return make([]string, 0)
}

func WriteFileShops(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, _ := os.Create(filename)
  defer file.Close()
  encoder := gob.NewEncoder(file)
  encoder.Encode(memShop)

  fileIndex, _ := os.Create(filename+".index")
  defer fileIndex.Close()
  encoderIndex := gob.NewEncoder(fileIndex)
  encoderIndex.Encode(memShopReg)
}


func LoadFileShops(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, err := os.Open(filename)
  if err !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename, err)
    return
  }
  defer file.Close()
  
  decoder := gob.NewDecoder(file)
  err = decoder.Decode(&memShop)
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
  errI = decoderI.Decode(&memShopReg)
  if errI != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename+".index", errI)
    return
  }
}

func GetShops(offset int, limit int) []Shop {
  res := make([]Shop, limit)
  i := 0
  for _, item := range memShop {
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

func ShopGQL() {
  AppendFields2GraphQL("shop", &graphql.Field{
			Type: shopType,
      Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                  Description: "id of the shops",
                  Type:graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id, ok := p.Args["id"].(string)
        if !ok {
          return nil, errors.New("Need ID")
        }
				return GetShopByID(id), nil
			},
		})
    
	AppendFields2GraphQL("shops", &graphql.Field{
			Type: graphql.NewList(shopType),
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
        glog.Infof("LOG: Shops(offset=%d, limit=%d)", offset, limit)
				return GetShops(offset, limit), nil
			},
    })
}

