$(function () {
  let $menubar = $('body>.main>.header')
  let $workspace = $('body>.main>.workspace')
  let menurows = (curr, dir) => {
    return (3 + Number(curr) + Number(dir)) % 3
  }

  $menubar
    .on('click', '>.menuitem', e => {
      let $t = $(e.currentTarget)

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
        .addClass('active')
        .trigger('activate')

      $t.addClass('selected')

      document.location.hash = $t.attr('name')
    })
    .on('click', '>.menu-scroll>div[dir]', e => {
      let $h = $(e.delegateTarget)

      $h.attr('menu-row', menurows($h.attr('menu-row'), $(e.currentTarget).attr('dir')))
        .children()
        .first()
        .click()
    })

  $(document.body).on('set', '.static.date', (e, d) => {
    $(e.currentTarget)
      .data('value', d)
      .text(d.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, ''))
  })

  $(document.body).on('reset', '.static.date', e => {
    $(e.currentTarget).text(
      $(e.currentTarget)
        .data('value')
        .replace('T', ' ')
        .replace(/:\d{1,2}(\..+)?Z.*/, ''))
  })

  $('body').on('error-message', (e, ...data) => {
    console.log('data', ...data)
  })
})
