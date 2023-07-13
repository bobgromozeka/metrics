package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/bobgromozeka/metrics/internal/helpers"
	"github.com/bobgromozeka/metrics/internal/metrics"

	"github.com/go-resty/resty/v2"
)

func reportToServer(serverAddr string, rm runtimeMetrics) {

	payloads := makeBodiesFromStructure(rm)

	if len(payloads) < 1 {
		return
	}

	client := resty.New()

	encodedPayload, err := json.Marshal(payloads)
	if err != nil {
		log.Println("Could not encode request: ", err)
		return
	}

	gzippedPayload, gzErr := helpers.Gzip(encodedPayload)
	if gzErr != nil {
		log.Println("Could not gzip request: ", gzErr)
		return
	}

	_, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(gzippedPayload).
		Post(serverAddr + "/updates")

}

func makeBodiesFromStructure(rm any) []metrics.RequestPayload {
	v := reflect.ValueOf(rm)
	t := reflect.TypeOf(rm)

	var payloads []metrics.RequestPayload

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fieldV := v.Field(i)
			fieldT := t.Field(i)
			if payload := makeBodyFromStructField(fieldV, fieldT); payload != nil {
				payloads = append(payloads, *payload)
			}
		}
	}

	return payloads
}

func makeBodyFromStructField(v reflect.Value, t reflect.StructField) *metrics.RequestPayload {
	metricsType := metrics.GaugeType
	if mt, ok := runtimeMetricsTypes[t.Name]; ok {
		metricsType = mt
	}

	rp := metrics.RequestPayload{
		ID:    t.Name,
		MType: metricsType,
	}

	//Shit conversions, but we lose accuracy anyway converting uint64 to float64
	switch metricsType {
	case metrics.GaugeType:
		switch val := v.Interface().(type) {
		case float64:
			rp.Value = &val
		case uint64, uint32:
			strVal := fmt.Sprintf("%d", v.Interface())
			intVal := helpers.StrToInt(strVal)
			fVal := float64(intVal)
			rp.Value = &fVal
		}
	case metrics.CounterType:
		strVal := fmt.Sprintf("%d", v.Interface())
		intVal := helpers.StrToInt(strVal)
		val := int64(intVal)
		rp.Delta = &val
	}

	if rp.Value == nil && rp.Delta == nil {
		return nil
	}

	return &rp
}
