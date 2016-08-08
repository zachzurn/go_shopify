package shopify


import (
  
    "bytes"
  
    "encoding/json"
  
    "fmt"
  
    "time"

    "io/ioutil"
  
)


type Asset struct {
  
    Attachment string `json:"attachment,omitempty"`

    ContentType string `json:"content_type,omitempty"`
  
    CreatedAt time.Time `json:"created_at,omitempty"`

    Key string `json:"key,omitempty"`
  
    Size int64 `json:"size,omitempty"`
  
    SourceKey string `json:"source_key,omitempty"`
  
    Src string `json:"src,omitempty"`

    ThemeId int64 `json:"theme_id,omitempty"`
  
    UpdatedAt time.Time `json:"updated_at,omitempty"`
  
    Value string `json:"value,omitempty"`
  

  
    api *API
  
}

type AssetUpload struct {
  
    Attachment string `json:"attachment,omitempty"`

    Key string `json:"key,omitempty"`
  
    Value string `json:"value,omitempty"`

  
    api *API
  
}



func (api *API) Assets(themeId int64) ([]Asset, error) {
  
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", themeId)
  res, status, err := api.request(endpoint, "GET", nil, nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := &map[string][]Asset{}
  err = json.NewDecoder(res).Decode(r)

  result := (*r)["assets"]

  if err != nil {
    return nil, err
  }

  for _, v := range result {
    v.api = api
  }

  return result, nil
}




func (api *API) Asset(themeId int64, assetKey string) (*Asset, error) {
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json?asset[key]=%s&theme_id=%d", themeId, assetKey, themeId )

  res, status, err := api.request(endpoint, "GET", nil, nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := map[string]Asset{}
  err = json.NewDecoder(res).Decode(&r)

  result := r["asset"]

  if err != nil {
    return nil, err
  }

  result.api = api

  return &result, nil
}


func (api *API) NewAsset() *Asset {
  return &Asset{api: api}
}

func (api *API) NewAssetUpload() *AssetUpload {
  return &AssetUpload{api: api}
}


func (obj *Asset) Save() (error) {
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", obj.ThemeId)
  method := "PUT"
  expectedStatus := 200

  body := map[string]*Asset{}
  body["asset"] = obj

  buf := &bytes.Buffer{}
  err := json.NewEncoder(buf).Encode(body)

  if b, err := ioutil.ReadAll(buf); err == nil {
      fmt.Printf(string(b))
  }

  if err != nil {
    return err
  }
  
  res, status, err := obj.api.request(endpoint, method, nil, buf)

  if err != nil {
    return err
  }

  if status != expectedStatus {
    r := errorResponse{}
    err = json.NewDecoder(res).Decode(&r)
    if err == nil {
      return fmt.Errorf("Status %d: %v", status, r.Errors)
    } else {
      return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
    }
  }

  r := map[string]Asset{}
  err = json.NewDecoder(res).Decode(&r)

  if err != nil {
    return err
  }

  *obj = r["asset"]

  return nil
}

func (obj *AssetUpload) Upload(themeId int64) (error) {
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", themeId)
  method := "PUT"
  expectedStatus := 200

  body := map[string]*AssetUpload{}
  body["asset"] = obj

  buf := &bytes.Buffer{}
  err := json.NewEncoder(buf).Encode(body)
  
  if err != nil {
    return err
  }
  
  res, status, err := obj.api.request(endpoint, method, nil, buf)

  if err != nil {
    return err
  }

  if status != expectedStatus {
    r := errorResponse{}
    err = json.NewDecoder(res).Decode(&r)
    if err == nil {
      return fmt.Errorf("Status %d: %v", status, r.Errors)
    } else {
      return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
    }
  }

  r := map[string]AssetUpload{}
  err = json.NewDecoder(res).Decode(&r)

  if err != nil {
    return err
  }

  *obj = r["asset"]

  return nil
}


func (api *API) Delete(themeId int64, assetKey string) (error) {
  endpoint := fmt.Sprintf("/admin/themes/%d/assets.json?asset[key]=%s", themeId, assetKey )

  res, status, err := api.request(endpoint, "DELETE", nil, nil)

  _ = res

  if err != nil {
    return err
  }

  if status == 403{
    return fmt.Errorf("Shopify does not allow this asset to be deleted.")    
  }

  if status != 200 {
    return fmt.Errorf("Status returned: %d", status)
  }

  return nil
}





