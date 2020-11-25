package main

import (
  "sync"
  "encoding/gob"
  "errors"
  "os"
  "unsafe"
  "github.com/golang/glog"
  "github.com/graphql-go/graphql"
  
)

type Region struct {
  ID           int64     `db:"id"             json:"id"              yaml:"id"`
  Name         string    `db:"title"          json:"title"          yaml:"title"`
}

var regionType = graphql.NewObject(
   graphql.ObjectConfig{
      Name: "Region",
      Fields: graphql.Fields{
         "id": &graphql.Field{
            Type: graphql.String,
         },
         "name": &graphql.Field{
            Type: graphql.String,
         },
      },
   },
)

var memRegion map[int64]Region

func RegionInit(max int64) {
  memRegion = make(map[int64]Region, max)
  glog.Infof("LOG: Region: max          = %d", max)
  glog.Infof("LOG: Region: sizeof(item) = %d", unsafe.Sizeof(Region{}))
  glog.Infof("LOG: Region: sizeof(map)  = %d", unsafe.Sizeof(memRegion))
}

func RegionCount() int64 {
  return int64(len(memRegion))
}

func RegionAppend(info *Region) {
  memRegion[info.ID] = *info
}

func GetRegionByID(id int64) (*Region) {
  item, ok := memRegion[id]
  if ok {
    glog.Infof("LOG: Region(%d) = %v", id, item)
    return &item
  }
  return nil
}

func GetRegions(offset int, limit int) []Region {
  res := make([]Region, limit)
  i := 0
  for _, item := range memRegion {
    if offset <= i &&  offset + limit > i {
      res = append(res, item)
    }
    if offset + limit < i {
      break
    }
    i++
  }
  glog.Infof("LOG: Regions %v", res)
  return res
}

func WriteFileRegions(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, _ := os.Create(filename)
  defer file.Close()
  encoder := gob.NewEncoder(file)
  encoder.Encode(memRegion)
}

func LoadFileRegions(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, err := os.Open(filename)
  if err !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename, err)
    return
  }
  defer file.Close()
  
  decoder := gob.NewDecoder(file)
  err = decoder.Decode(&memRegion)
  if err != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename, err)
    return
  }
}



func RegionGQL() {
  AppendFields2GraphQL("region", &graphql.Field{
			Type: regionType,
      Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                  Description: "id of the regions",
                  Type:graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id, ok := p.Args["id"].(int)
        if !ok {
          return nil, errors.New("Need ID")
        }
				return GetRegionByID(int64(id)), nil
			},
		})
    
	AppendFields2GraphQL("regions", &graphql.Field{
			Type: graphql.NewList(regionType),
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
        glog.Infof("LOG: Regions(offset=%d, limit=%d)", offset, limit)
				return GetRegions(offset, limit), nil
			},
    })
}
