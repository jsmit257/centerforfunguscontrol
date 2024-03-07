$(function () {
  var $menubar = $('body>.main>.header')
  var $workspace = $('body>.main>.workspace')

  $menubar.on('click', '>.menuitem', e => {
    var $t = $(e.currentTarget)
    $menubar
      .find('.menuitem')
      .removeClass('selected')
    $workspace
      .find('div')
      .removeClass('active')
    console.log('clicked', $workspace
      .find('>.' + $t.attr('name'))
      .trigger('activate'))
    $t.addClass('selected')
  })
})
