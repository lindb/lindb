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

* First line is no more than 50 characters and describes the changes.
* The body of the commit message should include a more detailed explanation of the change.
  It is ok to use markdown formatting in the explanation.

### Version Numbers

The LinDB version has three parts:

```
[major].[minor].[patch]
```