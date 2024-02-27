$(function () {
  var $menubar = $('body>.main>.header')
  var $workspace = $('body>.main>.workspace')

  $menubar.on('click', '.menuitem', e => {
    var $t = $(e.target)
    $menubar
      .find('.menuitem')
      .removeClass('selected')
    $workspace
      .find('div')
      .removeClass('active')
    $workspace
      .find('.' + $t.attr('name'))
      .addClass('active')
    $t.addClass('selected')
  })

  $workspace
    .find('.lifecycle>.ndx>.rows')
    .on('click', '.row', e => {
      var $row = $(e.currentTarget)
      console.log('clicked', $row.find('.id').text())
    })

  $workspace
    .find('.active>form>ul>li.expandable')
    .on('click', '', e => $(e.target).toggleClass('expanded'))

  // $workspace.find('.active>form>ul>li.expandable').toggleClass('expanded')

  console.log('asshole', $workspace.find('.active>form>ul>li.expandable'))
  // var $lifecycleForm = $workspace.find('div.lifecycle>form>ul.lifecycle')
  // $lifecycleForm.
  //   append('<li></li>').
  //   addClass("attr").
  //   append('<div></div>')
})
