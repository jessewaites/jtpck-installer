# Homebrew Tap Setup

## Steps to Complete Setup

### 1. Create Tap Repository
Create new GitHub repo: `jessewaites/homebrew-jtpck`
- Public repo
- No README, .gitignore, or license (GoReleaser will populate)

### 2. Create GitHub Token
1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token with `repo` scope
3. Copy token

### 3. Add Token to Secrets
In `jessewaites/jtpck-installer` repo:
1. Settings → Secrets and variables → Actions
2. New repository secret
3. Name: `HOMEBREW_TAP_TOKEN`
4. Value: [paste token]

### 4. Create First Release
```bash
git tag v0.1.0
git push origin v0.1.0
```

This triggers GitHub Actions:
- Builds binaries (darwin/linux, amd64/arm64)
- Creates GitHub release
- Updates homebrew-jtpck tap automatically

### 5. Test Installation
```bash
brew tap jessewaites/jtpck
brew install jtpck
jtpck --version
```

## Future Releases
```bash
git tag v0.2.0
git push origin v0.2.0
```

GoReleaser handles everything automatically.
