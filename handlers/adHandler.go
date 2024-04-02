package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/hunick1234/DcardBackend/config"
	"github.com/hunick1234/DcardBackend/dto"
	"github.com/hunick1234/DcardBackend/model/ad"
	"github.com/hunick1234/DcardBackend/model/ad/api/controller"
	v1 "github.com/hunick1234/DcardBackend/model/ad/api/v1"
	"github.com/hunick1234/DcardBackend/myhttp"
	"github.com/hunick1234/DcardBackend/server"
	"github.com/hunick1234/DcardBackend/service"
	"github.com/hunick1234/DcardBackend/types"
)

type AdHandlerFunc func(service.AdService, *types.AdControllerCtx) error

func AdHandler(server *server.Server, api AdHandlerFunc, flows []controller.APIController) http.HandlerFunc {
	conn, err := server.Pool.GetConnection(&config.MongoCfg{
		URI: "mongodb://localhost:27017",
		DB:  "dcard",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range flows {
		v.InitEvent(&types.AdControllerCtx{}, service.NewAdService(ad.NewAdRepository(conn)))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// create request context
		req, err := dto.NewRequest(r)
		if err != nil {
			log.Println("Error creating request:", err)
			httpError(w, http.StatusBadRequest, err.Error())
			return
		}
		adCtx := &types.AdControllerCtx{
			Ctx: r.Context(),
			R:   &req,
			W:   &myhttp.Response{},
			Err: nil,
		}

		conn, err := server.Pool.GetConnection(&config.MongoCfg{
			URI: "mongodb://localhost:27017",
			DB:  "dcard",
		})
		if err != nil {
			log.Println("Error getting connection:", err)
			httpError(w, http.StatusInternalServerError, "error getting connection")
			return
		}

		adRepo := ad.NewAdRepository(conn)
		adService := service.NewAdService(adRepo)

		start := time.Now()
		// before api event
		for _, flow := range flows {
			flow.BeforeAPIEvent(adCtx, adService)
		}
		elapsed := time.Since(start)
		log.Printf("bAPI took %s", elapsed)

		// api event
		start = time.Now()
		api(adService, adCtx)
		elapsed = time.Since(start)
		log.Printf("API took %s", elapsed)

		// after api event
		start = time.Now()
		for _, flow := range flows {
			flow.AfterAPIEvent(adCtx, adService)
		}
		elapsed = time.Since(start)
		log.Printf("aAPI took %s", elapsed)

		// write response
		if adCtx.Err != nil {
			httpError(w, adCtx.W.StausCode, adCtx.Err.Error())
			return
		} else {
			httpResponJson(w, adCtx.W.StausCode, adCtx.W.Message)
			return
		}
	}
}

func httpResponJson(w http.ResponseWriter, statusCode int, message []byte) {
	// write respon to json
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(message)
	w.Write(message)
	log.Println("<--", statusCode)
}

func httpError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println("Error writing response:", err)
	}
	log.Println("<--", statusCode)
}

func StartAPIControll(server *server.Server) {
	liveAd := controller.NewLiveAd()
	dailyAd := controller.NewDailyAd()

	server.HTTP.Get("/api/v1/ad", AdHandler(server, v1.GetAd, []controller.APIController{}))
	server.HTTP.Post("/api/v1/ad", AdHandler(server, v1.PostAd, []controller.APIController{dailyAd, liveAd}))
}
