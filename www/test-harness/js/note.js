$(_ => {
  let $rowtmpl = $('body>.main>.workspace>.template>.notes>.rows>.row.template')

  $('div.notes')
    .on('refresh', (e, owner) => {
      e.stopPropagation()

      let $n = $(e.currentTarget)

      $.ajax({
        url: `/notes/${owner}`,
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          $n
            .data('owner', owner)
            .find('>.rows')
            .trigger('send', result)
        },
        error: console.log,
        complete: (xhr, status) => {
          if ( // clumsy
            status === 'success'
            &&
            $n.find('.removable').length === 0
            &&
            $n.parents('.photos').length === 0
          ) {
            $n.find('div.add-note').click()
          }
        }
      })
    })
    .on('select', '>.rows', (e, id) => {
      let $rows = $(e.currentTarget)
        .find(`>.row[id=${id}]`)
        .click()

      if ($rows.find('.selected').length === 0) {
        $rows.children().first().click()
      }
    })
    .on('send', '>.rows', (e, ...data) => {
      e.stopPropagation()

      let $table = $(e.currentTarget)
      let selected = $table.find('>.selected').attr('id')

      $table.find('>.row.removable').remove()

      data.forEach(v => {
        $rowtmpl
          .clone(true, true)
          .insertBefore($table.children().first())
          .toggleClass('template removable')
          .data('row', v)
          .trigger('send', v)
      })

      $table.trigger('select', selected)
    })
    .on('send', '>.rows>.row', (e, data = { mtime: 'Now', ctime: 'Now' }) => {
      e.stopPropagation()

      let $row = $(e.currentTarget)

      $row.attr('id', data.id)
      $row.find('>.note').val(data.note)
      $row.find('>.mtime').trigger('set', data.mtime)
      $row.find('>.ctime').trigger('set', data.ctime)
    })
    .on('click', '>.rows>.row>div.action', e => {
      e.stopPropagation()

      let $n = $(e.currentTarget)
        .parents('.row')
        .first()
        .toggleClass('editing')

      if ($n.hasClass('adding')) {
        $n.remove()
      } else if (!$n.find('textarea').attr('disabled', !$n.hasClass('editing')).attr('disabled')) {
        $n.find('textarea').focus()
      }
    })
    .on('click', '>.rows>.row>div.commit', e => {
      e.stopPropagation()

      let $selected = $(e.currentTarget).parents('.row').first()
      let url = `/notes/${$(e.delegateTarget).data('owner')}`
      let other

      if ($selected.hasClass('adding')) {
        other = {
          method: 'POST',
          data: JSON.stringify({ note: $selected.find('>.note').val() }),
          error: _ => { $selected.remove() }
        }
      } else if ($selected.hasClass('editing')) {
        other = {
          method: 'PATCH',
          data: JSON.stringify({
            id: $selected.attr('id'),
            note: $selected.find('>.note').val(),
          }),
        }
      } else {
        other = {
          url: `${url}/${$selected.attr('id')}`,
          method: 'DELETE',
          success: (result, status, xhr) => {
            if ($selected.next().click().length === 0) {
              $selected.prev().click()
            }
            $selected.remove()
          },
        }
      }

      $.ajax({
        ...{
          url: url,
          method: 'HEAD',
          async: true,
          success: (result, status, xhr) => {
            $selected.trigger('send', result[0])
          },
          error: (xhr, status, err) => {
            $selected.trigger('reset')
            console.log(xhr, status, err)
          },
          complete: (xhr, status) => {
            $selected
              .removeClass('editing adding')
              .find('textarea')
              .attr('disabled', true)
          }
        },
        ...other
      })
    })
    .on('click', '>div.add-note', e => {
      e.stopPropagation()

      $rowtmpl
        .clone(true, true)
        .toggleClass('template removable editing adding')
        .insertBefore($(e.delegateTarget)
          .find('>.rows')
          .children()
          .first())
        .trigger('send')
        .find('textarea')
        .attr('disabled', false)
        .focus()
    })
})
