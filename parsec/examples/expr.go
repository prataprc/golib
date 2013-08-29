package main
import (
    "github.com/prataprc/golib/parsec"
    "os"
    "io/ioutil"
    "strconv"
    "fmt"
)

// Construct parser-combinator for parsing arithmetic expression on integer
func expr() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return ns[0]
        }
        return parsec.OrdChoice( "expr", nodify, false, sum )()(s)
    }
}

func prod() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodifyop := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return ns[0]
        }
        op := parsec.OrdChoice( "mdop", nodifyop, false, multop, divop )
        nodifyk := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return ns
        }
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns != nil {
                val := ns[0].(int)
                for _, x := range ns[1].([]parsec.ParsecNode) {
                    y := x.([]parsec.ParsecNode)
                    n := y[1].(int)
                    switch y[0].(*parsec.Terminal).Name {
                    case "MULT" : val *= n
                    case "DIV" : val /= n
                    }
                }
                return val
            }
            return nil
        }
        k := parsec.Kleene(
            "prod2kleene", nil,
            parsec.And("prod2", nodifyk, false, op, value),
        )
        return parsec.And( "prod", nodify, false, value, k )()(s)
    }
}
func sum() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodifyop := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return ns[0]
        }
        op := parsec.OrdChoice( "asop", nodifyop, false, addop, subop )
        nodifyk := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return ns
        }
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns != nil {
                val := ns[0].(int)
                for _, x := range ns[1].([]parsec.ParsecNode) {
                    y := x.([]parsec.ParsecNode)
                    n := y[1].(int)
                    switch y[0].(*parsec.Terminal).Name {
                    case "ADD" : val += n
                    case "SUB" : val -= n
                    }
                }
                return val
            }
            return nil
        }
        k := parsec.Kleene(
            "sum2kleene", nil,
            parsec.And("sum2", nodifyk, false, op, prod),
        )
        return parsec.And( "sub", nodify, false, prod, k )()(s)
    }
}

func groupExpr() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return ns[1]
        }
        return parsec.And(
            "groupExpr", nodify, false, openparan, expr, closeparan )()(s)
    }
}

func value() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil {
                return nil
            } else if n, ok := ns[0].(*parsec.Terminal); ok {
                if val, err := strconv.Atoi(n.Value); err == nil {
                    return val
                }
                fmt.Println("Invalid token", n.Tok)
                os.Exit(1)
            }
            return ns[0]
        }
        return parsec.OrdChoice(
            "value", nodify, false, parsec.Literal, groupExpr )()(s)
    }
}

var openparan = parsec.Terminalize( "(", "OPENPARAN", ")" )
var closeparan = parsec.Terminalize( ")", "CLOSEPARAN", ")" )
var addop = parsec.Terminalize( "+", "ADD", "+" )
var subop = parsec.Terminalize( "-", "SUB", "-" )
var multop = parsec.Terminalize( "*", "MULT", "*" )
var divop = parsec.Terminalize( "/", "DIV", "/" )

func main() {
    text, _ := ioutil.ReadFile(os.Args[1])
    s := parsec.NewGoScan(text)
    val := expr()(s)
    fmt.Println(val)
}
