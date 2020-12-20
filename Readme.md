# JUMPER

### How to add extra config to project json file
Add new project fields here `config/ProjectConfig.go` 
and add the related method to `config/config.go`, ie:
```
// SaveContainerNameToProjectConfig saves container name into project file
func (c *Config) SaveContainerNameToProjectConfig(cn string) (err error) {
	c.projectConfig.MainContainer = cn
	return c.fileSystem.SaveConfigFile(c.projectConfig, c.getProjectFile())
}

// GetProjectMainContainer gets project main container
func (c *Config) GetProjectMainContainer() string {
	return c.projectConfig.GetMainContainer()
}
```

### MIT License

#### Copyright (c) Anton Sannikov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.