package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	pb "github.com/je4/bremote/api"
	"github.com/je4/bremote/common"
	"github.com/mintance/go-uniqid"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func dummy(w http.ResponseWriter, r *http.Request) {
	return
}

func (controller *Controller) addRestRoutes(r *mux.Router) {

	// the proxy
	// ignore error because of static url, which must be correct
	proxy := &httputil.ReverseProxy{Director: controller.getProxyDirector()}
	proxy.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if controller.session == nil {
				return nil, errors.New("no tls session available")
			}
			return controller.session.Open()
		},
	}

	r.PathPrefix("/").HandlerFunc(common.MakePreflightHandler(
		controller.log,
	)).Methods("OPTIONS")

	r.HandleFunc("/", dummy)
	r.HandleFunc("/groups", controller.RestGroupList()).Methods("GET")
	r.HandleFunc("/groups/{group}", controller.RestGroupGetMember()).Methods("GET")
	r.HandleFunc("/groups/{group}", controller.RestGroupAddInstance()).Methods("PUT")
	r.HandleFunc("/groups/{group}", controller.RestGroupDelete()).Methods("DELETE")
	r.HandleFunc("/kvstore", controller.RestKVStoreExport()).Methods("GET")
	r.HandleFunc("/kvstore", controller.RestKVStoreImport()).Methods("POST", "PUT")
	r.HandleFunc("/kvstore/{client}", controller.RestKVStoreClientList()).Methods("GET")
	r.HandleFunc("/kvstore/{client}/{key}", controller.RestKVStoreClientValue()).Methods("GET")
	r.HandleFunc("/kvstore/{client}/{key}", controller.RestKVStoreClientValuePut()).Methods("PUT", "POST")
	r.HandleFunc("/kvstore/{client}/{key}", controller.RestKVStoreClientValueDelete()).Methods("DELETE")
	// get parameter: withstatus
	r.HandleFunc("/client", controller.RestClientList()).Methods("GET")
	r.HandleFunc("/client/{client}/status", controller.RestClientStatus()).Methods("GET")
	r.HandleFunc("/client/{client}/browserlog", controller.RestClientBrowserLog()).Methods("GET")
	r.HandleFunc("/client/{client}/addr", controller.RestClientHTTPSAddr()).Methods("GET")
	r.HandleFunc("/client/{client}/navigate", controller.RestClientNavigate()).Methods("POST")
	r.HandleFunc("/controller", controller.RestControllerList()).Methods("GET")
	r.HandleFunc("/controller/{controller}/templates", controller.RestControllerTemplates()).Methods("GET")
	r.PathPrefix("/{target}/").Handler(proxy).Methods("GET", "POST")

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.log.Infof(r.RequestURI)
			//			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})

	/*
	   	headersOk := handlers.AllowedHeaders([]string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Access-Control-Request-Method", "Authorization"})
	   	originsOk := handlers.AllowedOrigins([]string{"*"})
	   	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	   	credentialsOk := handlers.AllowCredentials()
	   //	ignoreOptions := handlers.IgnoreOptions()
	   	r.Use(handlers.CORS(
	   		originsOk,
	   		headersOk,
	   		methodsOk,
	   		credentialsOk,
	   //		ignoreOptions,
	   	))
	*/

}

func (controller *Controller) RestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.log.Infof(r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	}
}

func (controller *Controller) getProxyDirector() func(req *http.Request) {
	target, _ := url.Parse("http://localhost:80/")
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		//		vars := mux.Vars(req)
		//		t := vars["target"]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = common.SingleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set("X-Source-Instance", controller.GetInstance())
	}

	return director
}

func (controller *Controller) RestGroupList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestGroupList()")

		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		list, err := pw.GroupList(traceId)
		if err != nil {
			controller.log.Errorf("cannot get proxy group list: %v", err)
			http.Error(w, fmt.Sprintf("cannot get proxy group list: %v", err), http.StatusInternalServerError)
		}

		json, err := json.Marshal(list)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(json))
	}
}

func (controller *Controller) RestGroupGetMember() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestGroupGetMember()")
		vars := mux.Vars(r)
		group := vars["group"]

		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		list, err := pw.GroupGetMembers(traceId, group)
		if err != nil {
			controller.log.Errorf("cannot get members of group %v: %v", group, err)
			http.Error(w, fmt.Sprintf("cannot get members of group %v: %v", group, err), http.StatusInternalServerError)
		}

		json, err := json.Marshal(list)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(json))
	}
}

func (controller *Controller) RestGroupAddInstance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestGroupAddInstance()")
		vars := mux.Vars(r)
		group := vars["group"]

		decoder := json.NewDecoder(r.Body)
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			controller.log.Errorf("cannot decode data: %v", err)
			http.Error(w, fmt.Sprintf("cannot decode data: %v", err), http.StatusInternalServerError)
			return
		}
		d2, ok := data.(map[string]interface{})
		if !ok {
			controller.log.Errorf("invalid data format (not map[string]interface{})")
			http.Error(w, fmt.Sprintf("invalid data format (not map[string]interface{})"), http.StatusInternalServerError)
			return
		}
		d3, ok := d2["instance"]
		if !ok {
			controller.log.Errorf("invalid data format - no instance")
			http.Error(w, fmt.Sprintf("invalid data format - no instance"), http.StatusInternalServerError)
			return
		}
		instance, ok := d3.(string)
		if !ok {
			controller.log.Errorf("invalid data format - instance not a string")
			http.Error(w, fmt.Sprintf("invalid data format - instance not a string"), http.StatusInternalServerError)
			return
		}

		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		err = pw.GroupAddInstance(traceId, group, instance)
		if err != nil {
			controller.log.Errorf("cannot add instance %v to group %v: %v", instance, group, err)
			http.Error(w, fmt.Sprintf("cannot add instance %v to group %v: %v", instance, group, err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, `{"status":"ok"}`)
	}
}

func (controller *Controller) RestGroupDelete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestGroupAddInstance()")
		vars := mux.Vars(r)
		group := vars["group"]

		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		err := pw.GroupDelete(traceId, group)
		if err != nil {
			controller.log.Errorf("cannot delete group %v: %v", group, err)
			http.Error(w, fmt.Sprintf("cannot delete group %v: %v", group, err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, `"status":"ok"`)
	}
}

func (controller *Controller) RestClientHTTPSAddr() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestClientHTTPSAddr()")
		vars := mux.Vars(r)
		client, ok := vars["client"]
		if !ok {
			controller.log.Errorf("no client")
			http.Error(w, fmt.Sprintf("no client"), http.StatusInternalServerError)
		}

		cw := pb.NewClientWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		ret, err := cw.GetHTTPSAddr(traceId, client)
		if err != nil {
			controller.log.Errorf("cannot get status of %v: %v", client, err)
			http.Error(w, fmt.Sprintf("cannot get status of %v: %v", client, err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, ret)
	}
}

func (controller *Controller) RestClientStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestClientStatus()")
		vars := mux.Vars(r)
		client, ok := vars["client"]
		if !ok {
			controller.log.Errorf("no client")
			http.Error(w, fmt.Sprintf("no client"), http.StatusInternalServerError)
		}

		cw := pb.NewClientWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		ret, err := cw.GetStatus(traceId, client)
		if err != nil {
			controller.log.Errorf("cannot get status of %v: %v", client, err)
			http.Error(w, fmt.Sprintf("cannot get status of %v: %v", client, err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, ret)
	}
}

func (controller *Controller) RestClientBrowserLog() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestClientBrowserLog()")
		vars := mux.Vars(r)
		client, ok := vars["client"]
		if !ok {
			controller.log.Errorf("no client")
			http.Error(w, fmt.Sprintf("no client"), http.StatusInternalServerError)
		}

		cw := pb.NewClientWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		ret, err := cw.GetBrowserLog(traceId, client)
		if err != nil {
			controller.log.Errorf("cannot get status of %v: %v", client, err)
			http.Error(w, fmt.Sprintf("cannot get status of %v: %v", client, err), http.StatusInternalServerError)
		}
		jsonstr, err := json.Marshal(ret)
		if err != nil {
			controller.log.Errorf("cannot marshall %v: %v", ret, err)
			http.Error(w, fmt.Sprintf("cannot marshall %v: %v", ret, err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(jsonstr))
	}
}

func (controller *Controller) RestClientList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestClientList()")
		_, ok := r.URL.Query()["withstatus"]
		clients, err := controller.GetClients(ok)
		if err != nil {
			controller.log.Errorf("cannot get clients: %v", err)
			//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			http.Error(w, fmt.Sprintf("cannot get clients: %v", err), http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(clients)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(json))
	}
}

func (controller *Controller) RestControllerTemplates() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestControllerTemplates()")
		vars := mux.Vars(r)
		cname, ok := vars["controller"]
		if !ok {
			controller.log.Errorf("no controller in url")
			//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			http.Error(w, fmt.Sprintf("no controller in url"), http.StatusNotFound)
			return
		}
		var templates []string
		var err error
		if cname == controller.instance {
			templates, err = controller.GetTemplates()
			if err != nil {
				controller.log.Errorf("cannot get templates of %v: %v", cname, err)
				http.Error(w, fmt.Sprintf("cannot get templates of %v: %v", cname, err), http.StatusInternalServerError)
			}
		} else {
			cw := pb.NewControllerWrapper(controller.instance, controller.GetSessionPtr())
			traceId := uniqid.New(uniqid.Params{"traceid_", false})
			templates, err = cw.GetTemplates(traceId, cname)
			if err != nil {
				controller.log.Errorf("cannot get templates of %v: %v", cname, err)
				http.Error(w, fmt.Sprintf("cannot get templates of %v: %v", cname, err), http.StatusInternalServerError)
			}
		}
		jsonstr, err := json.Marshal(templates)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(jsonstr))
	}
}

func (controller *Controller) RestControllerList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestControllerList()")
		clients, err := controller.GetControllers()
		if err != nil {
			controller.log.Errorf("cannot get controllers: %v", err)
			//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			http.Error(w, fmt.Sprintf("cannot get controllers: %v", err), http.StatusInternalServerError)
			return
		}
		clients = append(clients, common.ClientInfo{
			InstanceName: controller.instance,
			Status:       "",
			Type:         common.SessionType_Controller,
		})
		json, err := json.Marshal(clients)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(json))
	}
}

func (controller *Controller) RestKVStoreExport() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestKVStoreExport()")

		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		value, err := pw.KVStoreList(traceId)
		if err != nil {
			controller.log.Errorf("cannot get value: %v", err)
			http.Error(w, fmt.Sprintf("cannot get value: %v", err), http.StatusInternalServerError)
		}

		result := map[string]interface{}{}
		for key, val := range *value {
			var d interface{}
			json.Unmarshal([]byte(val), &d)
			result[key] = d
		}

		jsonstr, err := json.Marshal(result)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/jsonstr")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(jsonstr))
	}
}

func (controller *Controller) RestKVStoreImport() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.log.Info("RestKVStoreImport()")

		// encode/decode to check for valid json
		decoder := json.NewDecoder(r.Body)
		data := map[string]interface{}{}
		err := decoder.Decode(&data)
		if err != nil {
			controller.log.Errorf("cannot decode data: %v", err)
			http.Error(w, fmt.Sprintf("cannot decode data: %v", err), http.StatusInternalServerError)
			return
		}
		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		for key, value := range data {
			keys := strings.SplitN(key, "-", 2)
			if len(keys) != 2 {
				controller.log.Errorf("key %v has no '-': %v", key, err)
				http.Error(w, fmt.Sprintf("key %v has no '-': %v", key, err), http.StatusInternalServerError)
				return
			}
			jsonstr, err := json.Marshal(value)
			if err != nil {
				controller.log.Errorf("cannot marshal value %v: %v", value, err)
				http.Error(w, fmt.Sprintf("cannot marshal value %v: %v", value, err), http.StatusInternalServerError)
				return
			}
			if err := pw.KVStoreSetValue(traceId, keys[0], keys[1], string(jsonstr)); err != nil {
				controller.log.Errorf("cannot set value for key %v: %v", key, err)
				http.Error(w, fmt.Sprintf("cannot set value for key %v: %v", key, err), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, `{"status":"ok"}`)
	}
}

func (controller *Controller) RestKVStoreClientList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		client := vars["client"]
		controller.log.Infof("RestKVStoreClientList(%v)", client)

		pw := pb.NewProxyWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		value, err := pw.KVStoreClientList(client, traceId)
		if err != nil {
			controller.log.Errorf("cannot get value: %v", err)
			http.Error(w, fmt.Sprintf("cannot get value: %v", err), http.StatusInternalServerError)
		}

		result := map[string]interface{}{}
		for key, val := range *value {
			var d interface{}
			json.Unmarshal([]byte(val), &d)
			result[key] = d
		}

		json, err := json.Marshal(result)
		if err != nil {
			controller.log.Errorf("cannot marshal result: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal result: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(json))
	}
}

func (controller *Controller) RestKVStoreClientValue() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		client := vars["client"]
		key := vars["key"]
		controller.log.Infof("RestKVStoreClientValue(%v, %v)", client, key)

		value, err := controller.GetVar(client, key)
		if err != nil {
			controller.log.Errorf("cannot get value: %v", err)
			http.Error(w, fmt.Sprintf("cannot get value: %v", err), http.StatusInternalServerError)
		}

		json, err := json.Marshal(value)
		if err != nil {
			controller.log.Errorf("cannot marshal value: %v", err)
			http.Error(w, fmt.Sprintf("cannot marshal value: %v", err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, string(json))
	}
}

func (controller *Controller) RestKVStoreClientValuePut() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		client := vars["client"]
		key := vars["key"]
		controller.log.Debugf("RestKVStoreClientValuePost(%v, %v)", client, key)

		// encode/decode to check for valid json
		decoder := json.NewDecoder(r.Body)
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			controller.log.Errorf("cannot decode data: %v", err)
			http.Error(w, fmt.Sprintf("cannot decode data: %v", err), http.StatusInternalServerError)
			return
		}

		err = controller.SetVar(client, key, data)
		if err != nil {
			controller.log.Errorf("cannot set value: %v", err)
			http.Error(w, fmt.Sprintf("cannot set value: %v", err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, `{"status":"ok"}`)
	}
}

func (controller *Controller) RestKVStoreClientValueDelete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		client := vars["client"]
		key := vars["key"]
		controller.log.Infof("RestKVStoreClientValueDelete(%v, %v)", client, key)

		err := controller.DeleteVar(client, key)
		if err != nil {
			controller.log.Errorf("cannot delete value: %v", err)
			http.Error(w, fmt.Sprintf("cannot delete value: %v", err), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, `{"status":"ok"}`)
	}
}

type controllerClientNavigate struct {
	Url         string `json:"url"`
	Nextstatus  string `json:"nextstatus,omitempty"`
	Waitfor     string `json:"waitfor,omitempty"`
	Waittimeout string `json:waittimeout,omitempty`
	PosX        int64  `json:posx,omitempty`
	PosY        int64  `json:posy,omitempty`
	Element     string `json:element,omitempty`
}

func (controller *Controller) RestClientNavigate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		client := vars["client"]

		decoder := json.NewDecoder(r.Body)
		var data controllerClientNavigate
		err := decoder.Decode(&data)
		if err != nil {
			controller.log.Errorf("cannot decode data: %v", err)
			http.Error(w, fmt.Sprintf("cannot decode data: %v", err), http.StatusInternalServerError)
			return
		}
		u, err := url.Parse(data.Url)
		if err != nil {
			controller.log.Errorf("cannot parse url %v: %v", data.Url, err)
			http.Error(w, fmt.Sprintf("cannot parse url %v: %v", data.Url, err), http.StatusInternalServerError)
		}

		controller.log.Infof("%v::RestClientNavigate(%v, %v)", client, u.String(), data.Nextstatus)

		cw := pb.NewClientWrapper(controller.instance, controller.GetSessionPtr())
		traceId := uniqid.New(uniqid.Params{"traceid_", false})
		err = cw.Navigate(traceId, client, u, data.Nextstatus)
		if err != nil {
			controller.log.Errorf("cannot navigate to %v: %v", u.String(), err)
			http.Error(w, fmt.Sprintf("cannot navigate to %v: %v", u.String(), err), http.StatusInternalServerError)
			return
		}
		waittimeout,err := time.ParseDuration(data.Waittimeout)
		if err != nil {
			controller.log.Errorf("invalid timeout format %v: %v", data.Waittimeout, err)
			http.Error(w, fmt.Sprintf("invalid timeout format %v: %v", data.Waittimeout, err), http.StatusInternalServerError)
			return
		}
		if int64(waittimeout) > 0 {
			if data.Element != "" {
				if err := cw.Click(traceId, client, data.Waitfor, waittimeout, data.Element); err != nil {
					controller.log.Errorf("cannot click: %v", err)
					http.Error(w, fmt.Sprintf("cannot click: %v", err), http.StatusInternalServerError)
					return
				}
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", controller.servername)
		io.WriteString(w, `{"status":"ok"}`)
	}
}
