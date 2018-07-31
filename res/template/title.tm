<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <title>Soro</title>
  <link href="https://fonts.googleapis.com/css?family=Exo+2:600" rel="stylesheet">
  <link href="/static/styles.css" rel="stylesheet">
</head>

<body>

  <div class="menu-outer">
    <span class="title">SORO</span>
    <div class="menu">
      <a href="/">ROOT</a>
    </div>
    <div class="menu">
      <a href="{{.dirPath}}">DIR ONLY</a>
    </div>
  </div>

  <div class="files-outer">
    {{.dir}}
  </div>

  <div class="thumbs-outer">
    {{.thumbs}}
  </div>

  <div class="{{.contentClass}}">
    <div class="preview-inner">
      <div class="preview-title">{{.contentTitle}}</div>
      <a href="{{.contentBack}}">
        <div class="preview-close">&lt;</div>
      </a>
      <div class="preview-content">{{.content}}</div>
    </div>
  </div>

</body>

</html>