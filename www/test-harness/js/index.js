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
})
