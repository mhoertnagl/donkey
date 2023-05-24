define i64 @main() {
main.entry:
        %0 = alloca i64
        store i64 0, i64* %0
        store i64 2, i64* %0
        %1 = load i64, i64* %0
        ret i64 %1
}