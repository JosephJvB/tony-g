/**
 * manually created an onChange trigger for the AppsScript in the AppsScript console thingo
 *
 * Alternatively I could sort these rows using the googlesheets go api
 * load all the rows
 * sort them
 * delete all rows
 * set new sorted rows
 * using the googlesheets go api I have
 */
function onChange(e) {
  var ss = SpreadsheetApp.getActiveSpreadsheet()
  var sheet = ss.getActiveSheet()

  switch (sheet.getName()) {
    case 'Apple Tracks':
      sortAppleTracks(sheet)
      break
    case 'Youtube Videos':
      sortYoutubeVideos(sheet)
      break
    case 'Youtube Tracks':
      sortYoutubeTracks(sheet)
      break
    case 'TEST':
      sortTest(sheet)
      break
  }
}

function sortAppleTracks(sheet) {
  var range = sheet.getRange('A2:F')
  range.sort([
    { column: 5, ascending: false }, // year
    { column: 6, ascending: false }, // added_at
  ])
}
function sortYoutubeVideos(sheet) {
  var range = sheet.getRange('A2:F')
  range.sort([
    { column: 3, ascending: false }, // published_at
    { column: 6, ascending: false }, // added_at
  ])
}
function sortYoutubeTracks(sheet) {
  var range = sheet.getRange('A2:I')
  range.sort([
    { column: 8, ascending: false }, // videopublish date
    { column: 9, ascending: false }, // added_at
  ])
}
function sortTest(sheet) {
  var range = sheet.getRange('A2:I')
  range.sort([
    { column: 8, ascending: false }, // videopublish date
    { column: 9, ascending: false }, // added_at
  ])
}
