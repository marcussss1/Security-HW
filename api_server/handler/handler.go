package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"security/parser"
	"security/store"
)

type Handler struct {
	store store.Store
}

func NewHandler(store store.Store) Handler {
	return Handler{store: store}
}

func (h *Handler) GetRequests(ctx echo.Context) error {
	requests := h.store.GetRequests()

	return ctx.JSON(http.StatusOK, requests)
}

func (h *Handler) GetRequestByID(ctx echo.Context) error {
	request := h.store.GetRequestByID(ctx.Param("id"))

	return ctx.JSON(http.StatusOK, request)
}

func (h *Handler) RepeatRequestByID(ctx echo.Context) error {
	request := h.store.GetRequestByID(ctx.Param("id"))
	repeatRequest := parser.ParseRepeatRequest(request)

	//go h.store.SaveRequest(repeatRequest)
	resp, err := http.DefaultTransport.RoundTrip(repeatRequest)
	if err != nil {
		return err
	}
	//go h.store.SaveResponse(resp)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) ScanRequestByID(ctx echo.Context) error {
	return nil
	//collection := h.Saver.GetClient().Database("test").Collection("requests")
	//
	//var result models.ReplyRequest
	//
	//params := mux.Vars(r)
	//ID, _ := params["id"]
	//
	//objectId, _ := primitive.ObjectIDFromHex(ID)
	//
	//filter := bson.D{{"_id", objectId}}
	//
	//_ = collection.FindOne(context.TODO(), filter).Decode(&result)
	//
	//Respond(w, 200, result)
}
