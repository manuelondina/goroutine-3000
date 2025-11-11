# Branch Protection Setup Instructions

This document provides step-by-step instructions for setting up branch protection rules on the `main` branch to ensure code quality and security.

## Overview

Branch protection rules prevent direct pushes to important branches and require pull requests to be reviewed before merging. This ensures that all code changes go through proper review and automated testing.

## Prerequisites

- Repository owner access (you must be `manuelondina` or have admin permissions)
- The GitHub Actions workflows must be set up in the repository (already done ✅)

## Step-by-Step Instructions

### Step 1: Navigate to Branch Protection Settings

1. Go to your repository on GitHub: https://github.com/manuelondina/goroutine-3000
2. Click on **Settings** (in the top navigation bar)
3. In the left sidebar, click on **Branches** (under "Code and automation")

### Step 2: Add Branch Protection Rule

1. Under "Branch protection rules", click **Add rule** or **Add branch protection rule**
2. In the "Branch name pattern" field, enter: `main`

### Step 3: Configure Protection Settings

Enable the following protection rules by checking these boxes:

#### Required Pull Request Reviews

- ✅ **Require a pull request before merging**
  - This ensures no one can push directly to main
  
- ✅ **Require approvals**
  - Set "Required number of approvals before merging" to: **1**
  - This means at least one approval is needed before merging
  
- ✅ **Dismiss stale pull request approvals when new commits are pushed**
  - This ensures reviews are for the latest code
  
- (Optional) ✅ **Require review from Code Owners**
  - Only if you create a CODEOWNERS file and want to enforce it

#### Designated Approvers

Since you (manuelondina) are the repository owner, you are automatically an eligible approver. To restrict who can approve:

1. ✅ **Restrict who can dismiss pull request reviews** (optional)
   - Add yourself or specific team members if you want fine-grained control

#### Required Status Checks

- ✅ **Require status checks to pass before merging**
  - This ensures CI/CD workflows pass before merging
  
- ✅ **Require branches to be up to date before merging**
  - This ensures the branch has the latest main branch changes
  
- In the search box that appears, search for and select:
  - **Run Tests with Coverage** (from the Go Test and Coverage workflow)
  - This makes the test workflow mandatory before merging

#### Additional Protections

- ✅ **Require conversation resolution before merging**
  - Ensures all PR comments are addressed
  
- ✅ **Require signed commits** (optional, but recommended for security)
  - Ensures all commits are cryptographically signed
  
- ✅ **Require linear history** (optional)
  - Prevents merge commits, requires rebase or squash merging
  
- ✅ **Include administrators**
  - This applies the rules even to repository administrators (including you!)
  - Highly recommended to ensure consistent practices

#### Restrict Push Access

- ✅ **Restrict who can push to matching branches**
  - Leave this empty or add specific users/teams who should have push access
  - With "Require a pull request before merging" enabled, direct pushes are blocked anyway

### Step 4: Save Protection Rule

1. Scroll to the bottom of the page
2. Click **Create** or **Save changes**

## Verification

After setting up the branch protection:

1. Try to push directly to main - it should be blocked:
   ```bash
   git checkout main
   git commit --allow-empty -m "Test protection"
   git push origin main
   # This should fail with a protection error
   ```

2. Create a pull request instead:
   ```bash
   git checkout -b test-branch
   git commit --allow-empty -m "Test PR process"
   git push origin test-branch
   # Then create PR on GitHub - it should require approval and passing tests
   ```

## Workflow Integration

With these protections in place:

1. **Go Test and Coverage workflow** will run automatically on every PR to main
2. The workflow must pass (all tests green) before merging is allowed
3. You (manuelondina) must approve the PR before it can be merged
4. Direct pushes to main are blocked - all changes must go through PRs

## Quick Reference: What's Protected

| Action | Allowed? | Requirements |
|--------|----------|--------------|
| Direct push to main | ❌ No | Must use PR |
| Merge PR without approval | ❌ No | Need 1+ approval |
| Merge PR with failing tests | ❌ No | Tests must pass |
| Merge PR without review | ❌ No | Review required |
| Create PR to main | ✅ Yes | Anyone with access |
| Review and approve PR | ✅ Yes | Repository collaborators |

## Additional Recommendations

1. **Create a develop branch** for ongoing development work
   ```bash
   git checkout -b develop
   git push origin develop
   ```

2. **Set up similar protection for develop** (optional)
   - Follow the same steps but with "develop" as the branch pattern
   - You might want less strict rules (e.g., no approval required)

3. **Set up CODEOWNERS file** (optional)
   - Create `.github/CODEOWNERS` to automatically request reviews from specific people
   - Example content:
     ```
     # Default owner for everything
     * @manuelondina
     
     # Specific owners for workflows
     .github/workflows/ @manuelondina
     ```

4. **Enable security features**
   - Go to Settings > Security > Code security and analysis
   - Enable Dependabot alerts and security updates
   - Enable Code scanning (GitHub Advanced Security)

## Troubleshooting

### "I can't merge even though I'm the owner"

- Check if "Include administrators" is enabled - this applies rules to everyone
- You can temporarily bypass by disabling protections, but this is not recommended

### "Status checks don't appear in the list"

- Status checks only appear after they've run at least once
- Create a test PR first to trigger the workflows
- Then return to branch protection settings to add the checks

### "Someone else needs to approve my PR"

- If you're the only collaborator and "Include administrators" is checked, you may need to:
  - Add another collaborator, OR
  - Temporarily disable "Require approvals" for your own PRs, OR
  - Use a GitHub App that can auto-approve (not recommended)

## Support

If you encounter issues with branch protection setup:
- Check GitHub's documentation: https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches
- Review GitHub's branch protection rules guide
- Contact GitHub Support if you believe there's a platform issue

---

**Note**: These instructions are specific to the `goroutine-3000` repository and assume you (manuelondina) are the repository owner with appropriate permissions.
