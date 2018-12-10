Error inspection [proposal](https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md) has two functions - `Is` and `As`. Both iterate over a chain of errors, both attempt to find an error in that chain. The only difference in those functions is that one works with values and another with types or so it seems.

**Problem**

When you use sentinel errors and try to wrap an error with a sentinel error its impossible to implement it in a way that allows `Is` to find that sentinel error in the error chain. You can see the problem [here](https://github.com/sidh/go2errorsinspection/blob/master/errors_test.go#L23). The problem stems from the fact that `Is` as currently proposed interacts with errors only in two ways - equality comparison and `Unwrap` method from `Wrapper` interface. To find an error in the chain you need both target error and error in the chain to be equal. To wrap an error you need to be able to actually store wrapped error somewhere. And you can't do both - you are either equal to an error that doesn't know anything about what it wrapping or you store an error that you wrap and stop being equal to the sentinel error you attempt to find.

With `As` and using typed errors you can easily test for errors wrapping another error like [here](https://github.com/sidh/go2errorsinspection/blob/master/errors_test.go#L51). 

`Is` behaviour is not a problem by itself - it might be exactly the right thing to do. The problem is that `Is` and `As` behave differently while being so similar on paper - their names are very alike, arguments are quite similar, both do about the same thing and have very similar documentation.

When you read the docs nothing suggests that their behaviour is different, you expect that they do exactly the same thing on exactly the same data (chain of errors), one works with types and another with values.

**Solutions**

There are several ways to solve the issue:
1. Remove `Is` altogether.
2. Rename `Is` to something else (`Cause`, whatever) and state specifically that it cannot find sentinel errors that wrap another error and that such behaviour is unwanted.
3. Allow errors to implement optional `Is` method.

**Opinion**

I don't think option 1 is viable as you have legitimate reasons for `Is` to exist and all it will accomplish is force people to implement their own `Is`.

Option 2 solves the issue nicely. It is especially good if you consider that in general typed errors are better than sentinel errors (you cannot overwrite typed error like a sentinel error, you can add behaviour to typed errors, etc). This option pushes people to use typed errors instead of sentinel errors. The problem is it might end up with the same result as option 1 - people will implement their own `Is` method.

Option 3 brings both `Is` and `As` on par in terms of functionality and behaviour but complicates things as is written in the original proposal.

Whatever the decision will be, if behaviour of `Is` is left as it is I strongly believe it should be better differentiated from `As` as to not cause confusion.
