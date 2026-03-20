1. **Add `Absolute` method to `EditList`**:
   - In `vidclip/edl.go`, I will add a new method `func (e EditList) Absolute(baseDir string) EditList` to `EditList`.
   - The method will create a new `EditList` slice with the same capacity as `e`.
   - It will iterate over `e`. For each `EditEntry`, it will check if `entry.Source` is an absolute path using `filepath.IsAbs()`.
   - If `entry.Source` is not absolute, it will use `filepath.Join(baseDir, entry.Source)` to make it absolute.
   - It will return the new `EditList`.

2. **Add tests for `Absolute` method**:
   - In `vidclip/edl_test.go`, I will add a new test function `TestAbsolute(t *testing.T)`.
   - The test will verify that relative paths are correctly prepended with `baseDir`.
   - The test will verify that absolute paths (e.g., `/already/absolute/path.mkv` or `C:\absolute\path.mkv` on Windows) remain unchanged.
   - It will also check that the original `EditList` remains unmodified (since the method returns a copy).

3. **Complete pre-commit steps**:
   - Follow `pre_commit_instructions` to ensure proper testing, verifications, reviews and reflections are done.

4. **Submit the change**.
