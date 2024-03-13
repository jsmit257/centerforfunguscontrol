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

  // $('.static.date').on('set', (e, d) => {
  //   console.log('this 1', e.currentTarget, 'date', d)
  //   // $(e.currentTarget).text(d.replace('T', ' ').replace(/(\.\d+)?Z/, ''))
  // })

  $(document.body).on('set', '.static.date', (e, d) => {
    // console.log('this 2', e.currentTarget, 'date', d)
    $(e.currentTarget)
      .text(d.replace('T', ' ').replace(/:\d{2}(\.\d+)?Z/, ''))
      .data('value', d)
  })

  $(document.body).on('reset', '.static.date', e => {
    // console.log('reset this', e.currentTarget, 'date', $(e.currentTarget)
    //   .data('value'))
    $(e.currentTarget).text(
      $(e.currentTarget)
        .data('value')
        .replace('T', ' ')
        .replace(/:\d{2}(\.\d+)?Z/, ''))
  })

  $('body').on('error-message', (e, ...data) => {
    console.log(data, ...data)
  })

})
