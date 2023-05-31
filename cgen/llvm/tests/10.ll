define i64 @main() {
main.entry:
        %0 = alloca i64
        store i64 0, i64* %0
        %1 = icmp slt i64 1, 2
        br i1 %1, label %if.then, label %if.merge

if.then:
        %2 = load i64, i64* %0
        %3 = add i64 %2, 1
        store i64 %3, i64* %0
        br label %if.merge

if.merge:
        %4 = icmp slt i64 2, 3
        br i1 %4, label %if.then.0, label %if.merge.0

if.then.0:
        %5 = load i64, i64* %0
        %6 = add i64 %5, 4
        store i64 %6, i64* %0
        br label %if.merge.0

if.merge.0:
        %7 = load i64, i64* %0
        ret i64 %7
}