package shopify

import (
  
    "encoding/json"
  
    "fmt"
  
    "time"
  
)

type Theme struct {
  
    CreatedAt time.Time `json:created_at`
  
    Id int64 `json:id`
  
    Name string `json:name`
  
    Role string `json:role`
  
    UpdatedAt time.Time `json:updated_at`
  
    Previewable string `json:previewable`
  
    Processing string `json:processing`
  
}


func (api *API) Theme_index() (*[]Theme, error) {
  res, status, err := api.request("/admin/themes.json", "GET", nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := &map[string][]Theme{}
  err = json.NewDecoder(res).Decode(r)

  fmt.Printf("things are: %v\n\n", *r)

  result := (*r)["products"]

	if err != nil {
		return nil, err
  }

  return &result, nil
}


// TODO implement Theme.show

// TODO implement Theme.create

// TODO implement Theme.update

// TODO implement Theme.destroy

