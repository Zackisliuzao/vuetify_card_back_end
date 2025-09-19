package models

type Dessert struct {
	Name     string  `json:"name"`
	Calories int     `json:"calories"`
	Fat      float64 `json:"fat"`
	Carbs    int     `json:"carbs"`
	Protein  float64 `json:"protein"`
	Iron     int     `json:"iron"`
}
