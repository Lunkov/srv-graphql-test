package main

import (
  "encoding/json"
  "github.com/graphql-go/graphql"
  "github.com/golang/glog"

)

var fieldsGQL = make(graphql.Fields)

func AppendFields2GraphQL(index string, f *graphql.Field) {
  fieldsGQL[index] = f
}

// https://medium.com/tunaiku-tech/what-is-graphql-and-how-is-it-implemented-in-golang-b2e7649529f1
// https://spec.graphql.org/June2018/
// https://blog.logrocket.com/3-tips-for-implementing-graphql-in-golang/
// https://habr.com/ru/company/ruvds/blog/444346/
// https://medium.com/tunaiku-tech/what-is-graphql-and-how-is-it-implemented-in-golang-b2e7649529f1

func funcGraphQL(query_str string) []byte  {
  glog.Infof("LOG: Query: %s", query_str)
  
  rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fieldsGQL}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		glog.Errorf("failed to create new schema, error: %v", err)
	}

	paramsGQL := graphql.Params{Schema: schema, RequestString: query_str}
	res := graphql.Do(paramsGQL)
	if len(res.Errors) > 0 {
		glog.Errorf("failed to execute graphql operation, errors: %+v", res.Errors)
	}
	rJSON, _ := json.Marshal(res)  
  return rJSON
}
