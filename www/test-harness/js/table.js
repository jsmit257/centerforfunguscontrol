$(function () {
  // var comparator = (insert, k, l, r) => {
  //   return (k, $l, $r) => {
  //     if (typeof $r === 'undefined') {
  //       return true
  //     }
  //     switch ($l.find(`.static.${k}`).text().localeCompare($r.find(`.static.${k}`).text())) {
  //       case insert:
  //       case 0:
  //         $l.insertBefore($r)
  //         return true
  //       default:
  //         return false
  //     }
  //   }
  // }

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
          args.buttonbar.find('.remove, .edit')[$table.children().length === 0 ? "removeClass" : "addClass"]('active')
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
        .addClass('selected editing')
      var $selected = $table
        .find('.selected')
        .removeClass('selected editing')
      if ($selected.length === 0) {
        $table.append($row)
      } else {
        $row.insertBefore($selected)
      }
      $row.find('input, select')
        .first()
        .focus()

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
              // async: typeof (args.async) === 'undefined' ? true : args.async,
              success: (result, status, xhr) => {
                args.success(result, status, xhr)
                if ($table.children().length > 0) {
                  args.buttonbar.find('.remove, .edit').addClass('active')
                }
              },
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
          "cancel": args.cancel || console.log,
          "ok": args.ok || (e => {
            var $selected = $table.find('.selected')
            $.ajax({
              url: args.url || `/${$table.attr('name')}/${$selected.find('>.uuid').text()}`,
              contentType: 'application/json',
              method: 'PATCH',
              dataType: 'json',
              data: args.data($selected),
              async: false,
              success: args.success,
              error: args.error || console.log,
            })
          })
        }
      })
    })
    .on('delete', (e, args) => {
      var $table = $(e.currentTarget)
      $.ajax({
        url: args.url || `/${$table.attr('name')}/${$table.find('.selected>.uuid').text()}`,
        contentType: 'application/json',
        method: 'DELETE',
        async: false,
        success: (result, status, xhr) => {
          $table.trigger('remove-selected')
          args.buttonbar.find('.remove, .edit')[$table.children().length === 0 ? "removeClass" : "addClass"]('active')
        },
        error: console.log,
      })
    })
    .on('sort', (e, key, order) => {
      var rows = []
      $(e.currentTarget)
        .children()
        .get()
        .forEach(v => { rows.push($(v)) })

      rows
        .sort(($l, $r) => {
          return (order === 'sort-desc' ? -1 : 1)
            * $l.find(`.static.${key}`).text()
              .localeCompare($r.find(`.static.${key}`).text())
        })
        .forEach(l => { $(e.currentTarget).append(l) })
    })
  // .data('sort-asc', comparator(-1))
  // .data('sort-desc', comparator(1))
  // .on('merge-sort', (e, key, order) => {
  //   var $table = $(e.currentTarget)
  //   var sort = $table.data(order || 'sort-asc')
  //   var rows = []
  //   $table.children().get().forEach(v => { rows.push($(v)) })
  //   var _ = function mergeSort(list) {
  //     if (list.length < 2) {
  //       return list
  //     }
  //     var middle = Math.floor(list.length / 2)
  //     var black = mergeSort(list.slice(0, middle))
  //     var red = mergeSort(list.slice(middle))
  //     for (var [i, j, before] = [0, red.length, 0]; i < j; i++) {
  //       while (!sort(key, red[i], black[before])) {
  //         before += 1
  //       }
  //       if (before === black.length) {
  //         red.slice(i).reverse().forEach(l => { l.insertAfter(black[before - 1]) })
  //         black.push(...red.slice(i))
  //         break
  //       }
  //       black.splice(before, 0, red[i])
  //     }
  //     return black
  //   }(rows)
  // })
})