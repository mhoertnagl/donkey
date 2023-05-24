define i64 @main() {
main.entry:
        %0 = alloca i64
        store i64 1, i64* %0
        %1 = alloca i64
        store i64 2, i64* %1
        %2 = load i64, i64* %1
        %3 = load i64, i64* %0
        %4 = icmp slt i64 %2, %3
        br i1 %4, label %if.then, label %if.merge

if.then:
        %5 = load i64, i64* %0
        ret i64 %5

if.merge:
        %6 = load i64, i64* %1
        ret i64 %6
}