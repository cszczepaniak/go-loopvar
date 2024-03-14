# go-loopvar
A static analysis tool to find unnecessary loop variable captures.

### Motivation
As of Go 1.22, it's no longer necessary to capture loop variables for use in closures, to take their
address, etc. Now loop variables are scoped to single iterations rather than the entire loop.

This tool finds cases that are probably unnecessary captures of loop variables. If run with `-fix`,
it will also remove the unnecessary captures.

### Example
This is the output of running `go-loopvar -fix` on a file with unnecessary captures.

```diff
diff --git a/test/something.go b/test/something.go
index dd5a7f9..67f179e 100644
--- a/test/something.go
+++ b/test/something.go
@@ -4,10 +4,8 @@ import "fmt"

 func foo() {
        for _, n := range []int{1, 2, 3} {
-               myNum := n
-
                go func() {
-                       fmt.Println(myNum)
+                       fmt.Println(n)
                }()
        }
 }
 ```
