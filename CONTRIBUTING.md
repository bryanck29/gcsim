# Contributing Guidelines

If you would like to contribute code to this project please follow these Pull Request guidelines:

0. (Optional but encouraged) Find at least one maintainer interested in your PR
1. Fork the project
2. Create a branch specifically for the feature you are contributing
3. (Optional but encouraged) Rebase your branch as needed. Please see quick reference if you are new to git.
4. After you are happy with your work, please make sure to submit a Pull Request from the feature branch. You are heavily discouraged from making pull requests from your main branch because we may not be able to get to your PR before you make new changes to your PR.


If contributing code for characters, please try to be explicit about what is complete or incomplete.
A rough guideline for your comment is as follows:

- [ ] Frames
- [ ] Talents
- [ ] C0 Normals
- [ ] C0 Skill
- [ ] C0 Burst
- [ ] A1
- [ ] A2 
- [ ] A4
- [ ] Sanity Test Cases
- [ ] C1
- [ ] C2
- [ ] C3
- [ ] C4
- [ ] C5
- [ ] C6

<details><summary>Click to expand copy-paste friendly version</summary>
  
```
- [ ] Frames
- [ ] Talents
- [ ] C0 Normals
- [ ] C0 Skill
- [ ] C0 Burst
- [ ] A1
- [ ] A2 
- [ ] A4
- [ ] Sanity Test Cases
- [ ] C1
- [ ] C2
- [ ] C3
- [ ] C4
- [ ] C5
- [ ] C6
```
</details>
Items may be omitted when irrelevant as is the case for many ascension passives.

# Git/Github Quick Reference Guide
For those who are new to git/github

```git checkout -b newbranchname``` creates a new branch from your current branch
If you have committed code, but upon finishing your feature, the main branch has progressed, you are encouraged to rebase it to ensure it still works. 
Please reach out for help if you are not sure how to do this step, the following steps can be dangerous and you can lose your work if not done correctly.

To rebase your branch you will need to run the command
```
git rebase --onto <newparent> <oldparent>
git push -f
```
Where new parent is the commitment hash of the newest commit on the main branch and old parent is the commitment hash of the oldest common commitment between your feature branch and the main branch.

# Local Testing/Building
0. Create a Config file
1. Navigate to ```./gcsim/cmd/gcsim```
2. Run ```go build``` to build the executable and then feed your config file in e.g. ```./gcsim.exe -c="config.txt" -out="out.gz" -gz``` OR Run ```go run . -c="config.txt" -out="out.gz" -gz``` 
3. Upload to the [viewer](https://viewer.gcsim.app) to view the output file to confirm everything is working accordingly, and optionally share the viewer in discord for debugging help.
