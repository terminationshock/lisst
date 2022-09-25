#!/bin/bash
set -e

mkdir -p test
rm -f test/*

export LISST_TEST=1

! ./lisst 2> /dev/null > test/RESULT
echo -e "foobar" | ./lisst --help > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 1 OK"

! echo -e "" | ./lisst 2> test/RESULT
echo "Empty input" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 2 OK"

! echo -e "\n\n" | ./lisst 2> test/RESULT
echo "Empty input" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 3 OK"

echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst > test/RESULT
echo -e "test1\n  test2\n    test3    test3" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 4 OK"

! echo "something" | ./lisst "no\match" 2> test/RESULT
echo "Invalid regular expression" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 5 OK"

echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst test > test/RESULT
echo -e "[red]test[-]1\n  [red]test[-]2\n    [red]test[-]3    test3" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 6 OK"

export LISST_COLOR=blue
echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst test > test/RESULT
echo -e "[blue]test[-]1\n  [blue]test[-]2\n    [blue]test[-]3    test3" > test/EXPECT
diff test/RESULT test/EXPECT
unset LISST_COLOR
echo "Test 7 OK"

echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst "t(est)" > test/RESULT
echo -e "t[red]est[-]1\n  t[red]est[-]2\n    t[red]est[-]3    test3" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 8 OK"

echo -e "test1\n  test2\n\n\ttest3\ttest3\n\n" | ./lisst nomatch > test/RESULT
echo -e "test1\n  test2\n    test3    test3" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 9 OK"

echo -e "#!/bin/sh\nexit 42" > test/EXIT
chmod +x test/EXIT
echo -e "something" | ./lisst some test/EXIT || echo "$?" > test/RESULT
echo "42" > test/EXPECT
diff test/RESULT test/EXPECT
rm test/EXIT
echo "Test 10 OK"

echo -e "test1 test2\ntest3" | ./lisst "test[1-9]" echo > test/RESULT
echo "test1" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 11 OK"

echo -e "test1 test2\ntest3" | ./lisst "test[1-9]" echo foobar > test/RESULT
echo "foobar test1" > test/EXPECT
diff test/RESULT test/EXPECT
echo "Test 12 OK"

rm -r test

echo "SUCCESS"
exit 0
