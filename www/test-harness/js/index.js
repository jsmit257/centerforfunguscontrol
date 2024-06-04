$(function () {
  let previous = {}
  let rownames = ['main', 'ancillary', 'reporting'/*,'analytics*/]
  let changerow = ((len) => {
    return (curr, dir) => { return rownames[(len + rownames.indexOf(curr) + Number(dir)) % len] }
  })(rownames.length)

  let $menubar = $('body>.main>.header')
    .on('click', '>.menuitem:not(.selected)', (e, data) => {
      let $t = $(e.currentTarget).addClass('selecting')
      let row = $(e.delegateTarget).attr('menu-row')
      let name = $t.attr('name')

      let $workspace = $(e.delegateTarget)
        .parent()
        .find('>.workspace')

      $menubar
        .find('>.menuitem.selected')
        .removeClass('selected')

      $workspace
        .find(`>div.active`)
        .removeClass('active')

      $workspace
        .find(`>.${$t.attr('tab-name') || name}`)
        .trigger('activate', data || $t)
        .addClass('active')

      $t.toggleClass('selected selecting')

      previous[row] = name

      document.location.hash = `${row}#${name}`
    })
    .on('select', '>.menu-scroll', (e, ...path) => {
      let data
      let btn
      switch (path.length) {
        case 3:
          data = path[2] // let it fall through
        case 2:
          if (!(btn = previous[path[0]])) {
            btn = previous[path[0]] = path[1]
          }
          $(e.delegateTarget).attr('menu-row', path[0])
          break
        case 1:
          // this really shouldn't be allowed: things go wrong if the
          // button isn't on the active row
          btn = path[0]
          break
      }

      $(e.delegateTarget).find(`>.menuitem[name="${btn}"]`).trigger('click', data)
    })
    .on('click', '>.menu-scroll>div[dir]', e => {
      let $h = $(e.delegateTarget)

      let next = changerow($h.attr('menu-row'), $(e.currentTarget).attr('dir'))
      $h.attr('menu-row', next)

      let btn = previous[next]
      if (!btn) {
        btn = previous[next] = $h.find('>.menuitem:visible').first().attr('name')
      }

      $h.find(`>[name="${btn}"]`).click()
    })

  $(document.body)
    .on('set', '.static.date', (e, d) => {
      $(e.currentTarget)
        .data('value', d)
        .text(d.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, ''))
    })
    .on('reset', '.static.date', e => {
      $(e.currentTarget).text(
        $(e.currentTarget)
          .data('value')
          .replace('T', ' ')
          .replace(/:\d{1,2}(\..+)?Z.*/, ''))
    })
    .on('error-message', (e, ...data) => {
      console.log('data', ...data)
    })
    .on('long-format', 'select.strain-list>option', (e, s) => {
      $(e.currentTarget).html(`${s.name} &bull; ${s.species} &bull; ${s.vendor.name} &bull; ${s.ctime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')}`)
    })
})
