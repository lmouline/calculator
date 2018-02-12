package resolver

import (
	"testing"
)

func BenchmarkResolver(b *testing.B) {
	for i := 0; i<b.N;i++ {
		Resolve("((1 + 1) - (4 * 5))*0 + 2 == (1+1.1)*2")
	}
}

func BenchmarkEvaluate(b *testing.B) {
	tests := []string{
		"1+1",
	//	//"π",
	//	"1+2^3^2",
	//	"2^(3+4)",
	//	"2^(3/(1+2))",
	//	"2^2(1+3)",
	//	"1+(-1)^2",
	//	"3*(3-(5+6)^12)*23^3-5^23",
	//	"2^3^2",
	//	//"ln(3^15)",
	//	//"sqrt(10)",
	//	//"abs(-3/2)",
	//	//"1+2sin(-1024)tan(acos(1))^2",
	//	//"tan(10)cos(20)",
	}
	for i := 0; i < b.N; i++ {
		Resolve(tests[i%len(tests)])
	}
}
