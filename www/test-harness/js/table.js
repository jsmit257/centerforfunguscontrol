$(function () {
  $('body>.main>.workspace>div .table>.rows')
    .on('click', '>.row', e => {
      var $row = $(e.currentTarget)
      if ($row.hasClass('selected')) {
        e.stopPropagation()
        return false
      }
      $row
        .parent()
        .find('.row.selected')
        .removeClass('selected editing')
      $row.addClass('selected')
    })
    .on('remove-selected', e => {
      var $selected = $(e.currentTarget).find('>.selected')
      console.log('selected', $selected)
      if ($selected.next().trigger('click').length == 0) {
        $selected.prev().trigger('click')
      }
      $selected.remove()
    })
    .on('refresh', (e, args) => {
      args ||= {}
      var $table = $(e.currentTarget)
      $.ajax({
        url: args.url || `/${$(e.currentTarget).attr('name')}s`,
        method: 'GET',
        async: true,
        success: (result, sc, xhr) => {
          var selected = $table.find('.selected>.uuid').text()
          $table.empty()
          result.forEach(r => {
            var $row = args.newRow(r)
            if (r.id === selected) {
              $row.addClass('selected')
            }
            $table.append($row)
          })
          if ($table.find('.selected').length == 0) {
            $table.find('.row').first().click()
          }
        },
        error: args.error || console.log
      })
    })
    .on('add', (e, args) => {
      var $table = $(e.currentTarget)

      var $row = args.newRow()
        .trigger('click')
        .addClass('editing')
        .find('input, select')
        .first()
        .focus()
      var $selected = $table
        .find('.selected')
        .removeClass('selected editing')
      if ($selected.length === 0) {
        $table.append($row)
      } else {
        $row.insertBefore($selected)
      }

      args.buttonbar.trigger('set', {
        "target": $table,
        "handlers": {
          "cancel": e => { $table.trigger('remove-selected') },
          "ok": args.ok || (e => {
            var $selected = $table.find('.selected')
            $.ajax({
              url: args.url || `/${$table.attr('name')}`,
              contentType: 'application/json',
              method: 'POST',
              dataType: 'json',
              data: args.data($selected),
              async: true,
              success: args.success,
              error: args.error,
            })
          })
        }
      })
    })
    .on('edit', (e, args) => {
      var $table = $(e.currentTarget)
      $table
        .find('.row.selected')
        .addClass('editing')
        .find('input, select')
        .first()
        .focus()

      args.buttonbar.trigger('set', {
        "target": $table,
        "handlers": {
          "cancel": console.log,
          "ok": args.ok || (e => {
            var $selected = $table.find('.selected')
            $.ajax({
              url: args.url || `/${$table.attr('name')}/${$selected.find('>.uuid').text()}`,
              contentType: 'application/json',
              method: 'PATCH',
              dataType: 'json',
              data: args.data($selected),
              async: true,
              success: args.success,
              error: args.error || console.log,
            })
          })
        }
      })
    })
    .on('delete', (e, url) => {
      var $table = $(e.currentTarget)
      $.ajax({
        url: url || `/${$table.attr('name')}/${$table.find('.selected>.uuid').text()}`,
        contentType: 'application/json',
        method: 'DELETE',
        async: true,
        success: (result, status, xhr) => {
          $table.trigger('remove-selected')
        },
        error: console.log,
      })
    })

})