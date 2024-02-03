package pkg

import(
	"log"
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("health checking")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "I am Amina and welcome to the place where I am gonna talk about sport teams of my fav comic(WindBreaker)🚴🏻‍♂️😇")
}

