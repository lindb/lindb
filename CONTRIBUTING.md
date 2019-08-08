## Contributing to LinDB

### Steps to Contribute

* Should you wish to work on an issue, please claim it first by commenting on the GitHub issue that you want to work on it. 
This is to prevent duplicated efforts from contributors on the same issue.
* Create a fork of the project(do not push branches directly to the main repository).
* Create a branch for your change.
* Make changes and add test(unit test/integration test).
* Commit the changes following the [commit guidelines](#git-commit-messages).
* Push the branch with your changes to your fork.
* Open a pull request against the LinDB project.

#### Testing

Where possible, test cases should be added to cover the new functionality or bug being
fixed. Test cases should be small, focused, and quick to execute.

#### Pull Requests

The following guidelines are to help ensure that pull requests (PRs) are easy to review and
comprehend.

* **One PR addresses one problem**, conflating issues in the same PR makes it more difficult
  to review and merge.
* **One commit per PR**, the final merge should have a single commit with a
  [good commit message](#git-commit-messages). Note, we can squash and merge via GitHub
  so it is fine to have many commits while working through the change and have us squash
  when it is complete. The exception is dependency updates where the
  only change is a dependency version. We typically do these as a batch with separate commits
  per version and merge without squashing. For this case, separate commits can be useful to
  allow for a git bisect to pinpoint a problem starting with a dependency change.
* **Reference related or fixed issues**, this helps us get more context for the change.
* **Partial work is welcome**, submit with a title including `[WIP]` (work in progress) to
  indicate it is not yet ready.
* **Keep us updated**, we will try our best to review and merge incoming PRs. We may close
  PRs after 30 days of inactivity. This covers cases like: failing tests, unresolved conflicts
  against master branch or unaddressed review comments.
  
  
### Issue Labels

For [issues][issue] we use the following labels to quickly categorize issues:

| Label Name     | Description                                                               |
|----------------|---------------------------------------------------------------------------|
| `bug`          | Confirmed bugs or reports that are very likely to be bugs.                |
| `enhancement`  | Feature requests.                                                         |
| `discussion`   | Requests for comment to figure out the direction.                         |
| `help wanted`  | Help from the community would be appreciated. Good first issues.          |
| `question`     | Questions more than bug reports or feature requests (e.g. how do I do X). |

### Git Commit Messages

Commit messages should try to follow these guidelines:

#### format

+ type: necessary(`:#issue` issue number is optional)
+ scope: optional
+ subject: necessary

```
[type:#issue][scope]: subject
```

#### type
 + feat: new feature
 + fix: bug solved
 + docs: commits of documentation
 + style: code style(The change does not affect code-logic)
 + refactor: neither new features nor bugfixs is added
 + test: changes of unit-test
 + chore: changes of ci-process and ci-tools

#### scope

scope is used to describe the scope of commit impact, such as `tsdb:index`, `broker:routes`, `storage`, etc

#### subject

subject is a short description of commit message, no more than 50 characters.

+ use simple present tense of the verbs, such as `change`, `add`, `fix`;
+ the first letter is lowercase;
+ end without a period(.);

#### Example
+ [feat]: add init commit
+ [chore:#1]: add travis ci
+ [test:#2]: add unit test for tsdb
+ [fix:#3][model:point]: change type timestamp from int64 to uint64
+ [feat:#4][tsdb]: implementation of memdb
+ [docs:#7]: add guideline for commit-message-format

### Version Numbers

The LinDB version has three parts:

```
[major].[minor].[patch]
```