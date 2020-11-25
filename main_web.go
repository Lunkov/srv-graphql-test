package main

import (
  "strconv"
  "net/http"
  "github.com/gorilla/mux"
  "encoding/json"
  "github.com/golang/glog"
)

func webGraphQL(w http.ResponseWriter, r *http.Request)  {
  keys, ok := r.URL.Query()["query"]
  if !ok || len(keys[0]) < 1 {
    glog.Infof("Url Param 'key' is missing")
    return
  }
       
  query_str := keys[0]  
  w.Write(funcGraphQL(query_str))
}

func webProduct4Sale(w http.ResponseWriter, r *http.Request)  {
  params := mux.Vars(r)
  
  
  product_id_str, okp1 := params["product_id"]
  if !okp1 {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  product_id, err1 := strconv.ParseInt(product_id_str, 10, 64)
  if err1 != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  
  region_id_str, okp2 := params["region_id"]
  if !okp2 {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  region_id, err2 := strconv.ParseInt(region_id_str, 10, 64)
  if err2 != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  
  res := GetProduct4Sale(product_id, region_id)
  
  jsonRes, _ := json.Marshal(res)
  w.Write(jsonRes)
}
