package main

var routeInfo []routeData

func main() {

	routeInf := buildRouteData(JUPPERI1)
	routeInfo = append(routeInfo, routeInf)

	routeInf = buildRouteData(JUPPERI2)
	routeInfo = append(routeInfo, routeInf)

	listenAndServe()
}
