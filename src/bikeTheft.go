package src

import (
    "time"
)

/*
    Class which models a bike theft
*/

type bikeTheft struct {
    Id  		int 		`json:"id"`
    Title 		string 		`json:"title"`
	Brand 		string		`json:"brand"`
    City 		string		`json:"city"`
    Description string		`json:"description"`
    Reported 	time.Time	`json:"reported"`
    Updated 	time.Time	`json:"updated"`
    Solved 		bool		`json:"solved"`
    OfficerId 	int			`json:"officer_id"`
    ImageName   string      `json:"image_name"`
    Image 		string		`json:"image"`
}
