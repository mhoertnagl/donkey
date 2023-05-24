define i64 @main() {
main.entry:
  %0 = call i64 @foo()
  ret i64 %0
}

define i64 @foo() {
foo.entry:
  ret i64 1
}