define i64 @main() {
main.entry:
        %0 = call i64 @foo(i64 3)
        ret i64 %0
}

define i64 @foo(i64 %a) {
foo.entry:
        %0 = alloca i64
        store i64 %a, i64* %0
        %1 = load i64, i64* %0
        %2 = add i64 %1, 1
        ret i64 %2
}