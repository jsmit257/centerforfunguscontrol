$(_ => {
  let newRow = (data = {}) => {
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="note static" />').text(data.id))
      .append($('<input class="note live">').val(data.note))
      .append($('<div class="mtime static const date" />')
        .data('value', data.mtime)
        .text(data.mtime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')))
      .append($('<div class="ctime static const date" />')
        .data('value', data.ctime)
        .text(data.ctime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')))
  }

  $('.notes.table')
    .on('init', $owner => {
      let $table = $(e.currentTarget)

      let $owned = $table.data('owner')
      if ($owned) {
        throw `this table is already owned by ${$owned}`
      }

      $table.data('owner', $owner)
    })
    .on('click', '>.row', e => {
      let $selected = $(e.currentTarget)

      $selected
        .parents('.table')
        .first()
        .find('>.selected')
        .removeClass('selected editing')

      $selected.addClass('selected')
    })
    .on('refresh', e => {
      let $table = $(e.currentTarget).empty()
      let selected = $table.find('>.selected>.uuid').text()

      $.ajax({
        url: `/notes/${$table.data('owner').find('.selected>.uuid').text()}`,
        method: 'GET',
        async: true,
        success: (result, status, xhr, $table = $(e.currentTarget).empty()) => {
          let $row
          result.foreach(v => { $table.append($row = newRow(v)) })
          if ($row.find('>.uuid').text() === selected) {
            $row.click()
          }
        },
        error: console.log,
      })
    })
    .on('add', e => {
      let $table = $(e.currentTarget)
      let $selected = $table.find('>.selected')

      $.ajax({
        url: `/notes/${$table.data('owner').find('.selected>.uuid').text()}`,
        method: 'POST',
        data: JSON.stringify({
          note: $selected.find('>.note.live').val(),
        }),
        async: true,
        success: (result, status, xhr) => {
          $selected.find('>.uuid').text(result[0].id)
          $selected.find('>.note.static').text(result[0].note)
          $selected.find('>.mtime.static').text(result[0].mtime)
          $selected.find('>.ctime.static').text(result[0].ctime)
        },
        error: console.log,
      })
    })
    .on('change', e => {
      let $table = $(e.currentTarget)
      let $selected = $table.find('>.selected')

      $.ajax({
        url: `/notes/${$table.data('owner').find('.selected>.uuid').text()}`,
        method: 'PATCH',
        data: JSON.stringify({
          id: $selected.find('>.uuid').text(),
          note: $selected.find('>.note.live').val(),
        }),
        async: true,
        success: (result, status, xhr) => {
          $selected.find('>.note.static').text(result[0].note)
          $selected.find('>.mtime.static').text(result[0].mtime)
        },
        error: console.log,
      })
    })
    .on('remove', e => {
      let $table = $(e.currentTarget)
      let $selected = $table.find('>.selected')

      $.ajax({
        url: `/notes/${$table.data('owner').find('.selected>.uuid').text()}/${selected.find('>.uuid').text()}`,
        method: 'DELETE',
        async: true,
        success: (result, status, xhr) => {
          if ($selected.nextSibling().click().length === 0) {
            $table.children().first().click()
          }
          $selected.remove()
        },
        error: console.log,
      })
    })
})
