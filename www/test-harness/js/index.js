$(function () {
  var $menubar = $('body>.main>.header')
  var $workspace = $('body>.main>.workspace')

  $menubar.on('click', '>.menuitem', e => {
    var $t = $(e.currentTarget)
    if ($t.hasClass('selected')) {
      return
    }
    $menubar
      .find('.menuitem')
      .removeClass('selected')
    $workspace
      .find('div')
      .removeClass('active')
    $workspace
      .find(`>.${$t.attr('name')}`)
      .trigger('activate')
    $t.addClass('selected')
  })

  $('body').on('error-message', (e, ...data) => {
    console.log(data, ...data)
  })

})
