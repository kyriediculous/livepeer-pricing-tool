package usecase

import (
	"encoding/json"
	"net/http"
	"time"

	"api/model"
	"api/dataservice"

	log "github.com/sirupsen/logrus"
)

// Specifies the API endpoint where data is exposed by livepeer broadcaster node
const broadcasterEndpoint = "http://localhost:7935/registeredOrchestrators"

// Specified the time duration (in seconds) between data polls.
const pollingInterval = 3600

// Fetches the data from broadcaster endpoint, and stores it in the pricing tool DB.
func GetData() []model.Orchestrator {
	log.Infoln("Fetching data from broadcaster endpoint.")
	response, err := http.Get(broadcasterEndpoint)
	if err != nil {
		log.Errorln("The HTTP request failed with error", err)
	} else {
		log.Infoln("The HTTP reqeust succeeded.")
	}

	orchestrators := []model.Orchestrator{}
	err = json.NewDecoder(response.Body).Decode(&orchestrators)
	if err != nil {
		log.Errorln("Error in JSON parsing", err.Error())
	} else {
		log.Infoln("JSON parsing successful.")
	}
	return orchestrators
}

// Adds orchestrator statistics and price history to the database 
func InsertInDB(orchestrators []model.Orchestrator) {

	for i, x := range orchestrators {
		if dataservice.IfOrchestratorExists(x.Address) {
			log.Infoln(i, "Updating orchestrator statistics for", x.Address)
			dataservice.UpdateOrchestrator(x)
		} else {
			log.Infoln(i, "Inserting orchestrator statistics for", x.Address)
			dataservice.InsertOrchestrator(x)
		}
		dataservice.InsertPriceHistory(x)	
	}
}

// Polls for data from the broadcaster endpoint at specified polling intervals
func PollForData() {
	log.Infoln("Polling service initiated.")
	for {
		InsertInDB(GetData())
		time.Sleep(pollingInterval * time.Second)
	}
}
