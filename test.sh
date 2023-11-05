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
        echo -e "[::r]test[::-]1\n  [::r]test[::-]2\n    [::r]test[::-]3    test3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    7)
        echo -e "[test]1\n  [test]2\n\n\ttest3\t[test]3\n\n" | ./lisst test > test/RESULT_$1
        echo -e "[[::r]test[::-][]1\n  [[::r]test[::-][]2\n    [::r]test[::-]3    [test[]3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    8)
        echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst "t(est)" > test/RESULT_$1
        echo -e "t[::r]est[::-]1\n  t[::r]est[::-]2\n    t[::r]est[::-]3    test3" > test/EXPECT_$1
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
        diff test/RESULT_$1 test/EXPECT_$1 && rm test/EXIT
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
    13)
        echo -e "test1 test2\ntest3" | ./lisst --line > test/RESULT_$1
        echo -e "[::r]test1 test2[::-]\n[::r]test3[::-]" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    14)
        echo -e "test1234 abcdef1\n0987654321 1234567\n--git-commit-hash" | ./lisst --git-commit-hash > test/RESULT_$1
        echo -e "test1234 [::r]abcdef1[::-]\n[::r]0987654321[::-] 1234567\n--git-commit-hash" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    15)
        echo -e "./src/main.go:content\nMakefile:42:content" | ./lisst --filename > test/RESULT_$1
        echo -e "[::r]./src/main.go[::-]:content\n[::r]Makefile[::-]:42:content" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    16)
        echo -e "./src/main.go:content\nMakefile:42:content" | ./lisst --filename echo > test/RESULT_$1
        echo -e "./src/main.go" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    17)
        ! echo -e "test1\nztest2" | ./lisst "test" invalid 2> test/RESULT_$1
        grep -q "file not found" test/RESULT_$1
        ;;
    18)
        echo -e "./src/main.go:content\nMakefile:42:content" | ./lisst --show-output --filename echo -n foo > test/RESULT_$1
        echo -e "foo ./src/main.go" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    19)
        echo -e "#!/bin/bash\necho \"\$@\";echo -n \"\$@\" 1>&2" > test/SCRIPT_$1
        chmod +x test/SCRIPT_$1
        echo -e "./src/main.go:content\nMakefile:42:content" | ./lisst --filename --show-output test/SCRIPT_$1 > test/RESULT_$1
        echo -e "./src/main.go\n./src/main.go" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    20)
        echo -e "test1 test2\ntest3" | ./lisst "test[1-9]" echo {} foobar > test/RESULT_$1
        echo "test1 foobar" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    21)
        echo -e "test1 test2\ntest3" | ./lisst "test[1-9]" echo foo{}b{}ar > test/RESULT_$1
        echo "footest1btest1ar" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    22)
        echo -e "extract the ./ file ./test.sh and not Makefile" | ./lisst --filename echo > test/RESULT_$1
        echo "./test.sh" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    23)
        echo -e "extract the ./ file ./test.sh and not Makefile" | ./lisst --dirname echo > test/RESULT_$1
        echo "./" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    24)
        echo -e "this text is \033[0;31mred\033[0m and \033[0;34mblue\033[0m" | ./lisst red > test/RESULT_$1
        echo "this text is [maroon:][::r]red[::-][-:-:] and [navy:]blue[-:-:]" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    25)
        echo -e "test1 test2\ntest3\ntest2" | ./lisst --filter test2 > test/RESULT_$1
        echo -e "test1 [::r]test2[::-]\n[::r]test2[::-]" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    26)
        echo -e "test1 test2\ntest3\ntest2" | ./lisst test2 --filter > test/RESULT_$1
        echo -e "test1 [::r]test2[::-]\n[::r]test2[::-]" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    27)
        ! echo -e "test1 test2\ntest3\ntest2" | ./lisst --filter no_match 2> test/RESULT_$1
        echo "All lines filtered out" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    28)
        echo -e "./src/main.go:content:2\nMakefile:42:content\nMakefile:1:2:3" | ./lisst --filename-lineno > test/RESULT_$1
        echo -e "./src/main.go:content:2\n[::r]Makefile:42[::-]:content\n[::r]Makefile:1[::-]:2:3" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    29)
        echo -e "test3\ntest11\nfoobar2\ntest2\nfoobar1\n" | ./lisst --sort "test[0-9]" > test/RESULT_$1
        echo -e "[::r]test1[::-]1\n[::r]test2[::-]\n[::r]test3[::-]\nfoobar2\nfoobar1" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    30)
        echo -e "test3\ntest11\nfoobar2\ntest2\nfoobar1\n" | ./lisst --sort-rev "test[0-9]" > test/RESULT_$1
        echo -e "[::r]test3[::-]\n[::r]test2[::-]\n[::r]test1[::-]1\nfoobar2\nfoobar1" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    31)
        echo -e "200\n3.14\n\n70.1\n" | ./lisst --sort --line > test/RESULT_$1
        echo -e "[::r]3.14[::-]\n[::r]70.1[::-]\n[::r]200[::-]" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    32)
        echo -e "foobar\n" | ./lisst --foobar 2> test/RESULT_$1
        echo "Invalid command-line option --foobar" > test/EXPECT_$1
        diff test/RESULT_$1 test/EXPECT_$1
        ;;
    esac
}

if [ $# -eq 0 ]; then
    result=0
    for i in {1..32}; do
        echo "Test $i"
        run $i || { result=1; echo "   FAILED"; }
    done
    if [ $result -eq 0 ]; then
        rm -r test/
    fi
    exit $result
fi
