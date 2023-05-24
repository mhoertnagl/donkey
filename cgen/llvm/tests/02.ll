define i64 @main() {
main.entry:
        %0 = alloca i64
        store i64 1, i64* %0
        %1 = alloca i64
        store i64 2, i64* %1
        %2 = load i64, i64* %0
        %3 = load i64, i64* %1
        %4 = add i64 %2, %3
        ret i64 %4
}