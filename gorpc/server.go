package gorpc

import (
    "context"
    "fmt"
    "github.com/baxiang/gorpc/interceptor"
    "log"

    "reflect"

)

type Server struct {
    opts *ServerOptions
    services map[string]Service
    closing bool
}


func NewServer(opt ...ServerOption)*Server{
    s :=&Server{
        opts:&ServerOptions{},
        services: map[string]Service{},
    }
    for _,o :=range opt{
        o(s.opts)
    }
    return s
}

type emptyInterface interface{}

func (s *Server)RegisterService(serviceName string, svr interface{}) error {
    svrType := reflect.TypeOf(svr)
    svrValue := reflect.ValueOf(svr)
    sd := &ServiceDesc{
        ServiceName: serviceName,
        // for compatibility with code generation
        HandlerType : (*emptyInterface)(nil),
        Svr : svr,
    }

    methods, err := getServiceMethods(svrType, svrValue)
    if err != nil {
        return err
    }
    sd.Methods = methods
    s.Register(sd, svr)

    return nil

}

func getServiceMethods(serviceType reflect.Type, serviceValue reflect.Value) ([]*MethodDesc, error) {

    var methods []*MethodDesc

    for i := 0; i < serviceType.NumMethod(); i++ {
        method := serviceType.Method(i)

        if err := checkMethod(method.Type); err != nil {
            return nil, err
        }

        methodHandler := func (svr interface{},ctx context.Context, dec func(interface{}) error, ceps []interceptor.ServerInterceptor) (interface{}, error) {

            reqType := method.Type.In(2)

            // determine type
            req := reflect.New(reqType.Elem()).Interface()

            if err := dec(req); err != nil {
                return nil, err
            }

            if len(ceps) == 0 {
                values := method.Func.Call([]reflect.Value{serviceValue,reflect.ValueOf(ctx),reflect.ValueOf(req)})
                // determine error
                return values[0].Interface(), nil
            }

            handler := func(ctx context.Context, reqbody interface{}) (interface{}, error) {

                values := method.Func.Call([]reflect.Value{serviceValue,reflect.ValueOf(ctx),reflect.ValueOf(req)})

                return values[0].Interface(), nil
            }

            return interceptor.ServerIntercept(ctx, req, ceps, handler)
        }

        methods = append(methods, &MethodDesc{
            MethodName: method.Name,
            Handler: methodHandler,
        })
    }

    return methods , nil
}

func checkMethod(method reflect.Type) error {

    // params num must >= 2 , needs to be combined with itself
    if method.NumIn() < 3 {
        return fmt.Errorf("method %s invalid, the number of params < 2", method.Name())
    }

    // return values nums must be 2
    if method.NumOut() != 2 {
        return fmt.Errorf("method %s invalid, the number of return values != 2", method.Name())
    }

    // the first parameter must be context
    ctxType := method.In(1)
    var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
    if !ctxType.Implements(contextType) {
        return fmt.Errorf("method %s invalid, first param is not context", method.Name())
    }

    // the second parameter type must be pointer
    argType := method.In(2)
    if argType.Kind() != reflect.Ptr {
        return fmt.Errorf("method %s invalid, req type is not a pointer", method.Name())
    }

    // the first return type must be a pointer
    replyType := method.Out(0)
    if replyType.Kind() != reflect.Ptr {
        return fmt.Errorf("method %s invalid, reply type is not a pointer", method.Name())
    }

    // The second return value must be an error
    errType := method.Out(1)
    var errorType = reflect.TypeOf((*error)(nil)).Elem()
    if !errType.Implements(errorType) {
        return fmt.Errorf("method %s invalid, returns %s , not error", method.Name(), errType.Name())
    }

    return nil
}

func (s *Server)Register(sd *ServiceDesc,svr interface{}){
    if sd == nil||svr ==nil{
        return
    }
    ht := reflect.TypeOf(sd.HandlerType)
    st := reflect.TypeOf(svr)
    if !st.Implements(ht){
        log.Fatal("handlerType %v not match service : %v ", ht, st)
    }
    ser :=&service{
        serviceName: sd.ServiceName,
        svr:         svr,
        handlers:    make(map[string]Handler),
    }
    for _, method := range sd.Methods {
        ser.handlers[method.MethodName] = method.Handler
    }
    s.services[sd.ServiceName] = ser
}


func (s *Server) Serve() {


}

type emptyService struct{}

func (s *Server) ServeHttp() {


}

func (s *Server) Close() {
    s.closing = false

    for _, service := range s.services {
        service.Close()
    }
}