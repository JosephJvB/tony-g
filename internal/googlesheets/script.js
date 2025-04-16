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

  if (sheet.getName() !== 'Scraped Tracks') {
    return
  }

  var range = sheet.getRange('A2:F')
  range.sort([
    { column: 4, ascending: false }, // year
    { column: 6, ascending: false }, // added_at
  ])
}
