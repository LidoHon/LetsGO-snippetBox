package main

import (
	"testing"
	"time"

	"github.com/LidoHon/LetsGO-snippetBox.git/internal/assert"
)

func TestHumanDate(t *testing.T){
	tests :=[]struct {
		name string
		tm time.Time
		want string
	}{
		{
			name: "UTC",
			tm : time.Date(2024, 11, 28, 10, 15, 0, 0, time.UTC),
			want: "28 Nov 2024 at 10:15",
		},
		{
			name:"Empty",
			tm: time.Time{},
			want:"",
		},
		{
			name: "CET",
			tm: time.Date(2024, 11, 28, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "28 Nov 2024 at 09:15",
		},

	}
	
	

	for _, tt :=range tests{
		t.Run(tt.name, func(t *testing.T){
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
		
	}
}