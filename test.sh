#!/bin/bash
set -e

mkdir -p test
export LISST_TEST=1

function run() {
    case $1 in
    1)
        ! ./lisst 2> /dev/null > test/RESULT_$1
        ./lisst --help > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1 || echo "   SKIPPED" # Works in interactive shells only
        ;;
    2)
        ! echo -e "" | ./lisst 2> test/RESULT_$1
        echo "Empty input" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    3)
        ! echo -e "\n\n" | ./lisst 2> test/RESULT_$1
        echo "Empty input" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    4)
        echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst > test/RESULT_$1
        echo -e "test1\n  test2\n    test3    test3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    5)
        ! echo "something" | ./lisst "no\match" 2> test/RESULT_$1
        echo "Invalid regular expression" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    6)
        echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst test > test/RESULT_$1
        echo -e "[red]test[-]1\n  [red]test[-]2\n    [red]test[-]3    test3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    7)
        export LISST_COLOR=blue
        echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst test > test/RESULT_$1
        echo -e "[blue]test[-]1\n  [blue]test[-]2\n    [blue]test[-]3    test3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        unset LISST_COLOR
        ;;
    8)
        echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst "t(est)" > test/RESULT_$1
        echo -e "t[red]est[-]1\n  t[red]est[-]2\n    t[red]est[-]3    test3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    9)
        echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst nomatch > test/RESULT_$1
        echo -e "test1\n  test2\n    test3    test3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    10)
        echo -e "#!/bin/sh\nexit 42" > test/EXIT
        chmod +x test/EXIT
        echo -e "something" | ./lisst some test/EXIT || echo "$?" > test/RESULT_$1
        echo "42" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        rm test/EXIT
        ;;
    11)
        echo -e "test1 test2\ntest3" | ./lisst "test[1-9]" echo > test/RESULT_$1
        echo "test1" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    12)
        echo -e "test1 test2\ntest3" | ./lisst "test[1-9]" echo foobar > test/RESULT_$1
        echo "foobar test1" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    esac
}

if [ $# -eq 0 ]; then
    result=0
    for i in {1..12}; do
        echo "Test $i"
        run $i || { result=1; echo "   FAILED"; }
    done
    if [ $result -eq 0 ]; then
        rm -r test/
    fi
    exit $result
fi
