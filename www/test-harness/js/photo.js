$(_ => {
  let $rowtmpl = $('body>.template>.photos>.rows>.row.template')

  $('body>.template>.notes')
    .clone(true, true)
    .appendTo($('body>.template>.photos>.rows>.row.template'))

  $('div.photos')
    .on('refresh', (e, owner) => {
      e.stopPropagation()

      let $n = $(e.currentTarget)

      $.ajax({
        url: `/photos/${owner}`,
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
          if (status === 'success' && $n.find('.removable').length === 0) {
            $n.find('div.add-photo').click()
          }
        }
      })
    })
    .on('send', '>.rows', (e, ...data) => {
      e.stopPropagation()

      let $table = $(e.currentTarget)

      $table.find('>.row.removable').remove()

      data.forEach(v => {
        $rowtmpl
          .clone(true, true)
          .insertBefore($table.children().first())
          .toggleClass('template removable')
          .data('row', v)
          .trigger('send', v)
          .find('>.photo-ctl')
          .addClass('active')
      })
    })
    .on('send', '>.rows>.row', (e, data = { mtime: 'Now', ctime: 'Now' }) => {
      e.stopPropagation()

      let $row = $(e.currentTarget)

      $row.attr('id', data.id)
      $row.find('>img.image').attr('src', `/album/${data.image || '../images/transparent.png'}`)
      $row.find('>.mtime').trigger('set', data.mtime)
      $row.find('>.ctime').trigger('set', data.ctime)
    })
    .on('click', '>.rows>.row>div.back', e => {
      let $row = $(e.currentTarget)
        .parents('.row')
        .first()
        .removeClass('selected noting')
        .parents('.photos')
        .toggleClass('gallery singleton')
    })
    .on('click', '>.rows>.row>div.toggle-notes', e => {
      let $row = $(e.currentTarget)
        .parents('.row')
        .first()
        .toggleClass('noting')

      if ($row.hasClass('noting')) {
        $row.find('>.notes').trigger('refresh', $row.attr('id'))
      }
    })
    .on('click', '>.rows>.row>div.paste-file', e => {
      $(e.currentTarget).trigger('paste')
    })
    .on('paste', (e, cd) => {
      $('.rows>.row>div.paste-file')
        .parents('.row')
        .first()
        .find('input.photo')
        .get(0)
        .files = cd.files;
    })
    .on('click', '>.rows>.row>div.action', e => {
      e.stopPropagation()

      let $n = $(e.currentTarget)
        .parents('.row')
        .first()
        .toggleClass('editing')

      if ($n.hasClass('adding')) {
        if (
          $n
            .parents('.rows')
            .first()
            .find('>.removable')
            .length === 0
        ) {
          $n
            .parents('.table.strain')
            .first()
            .removeClass('photoing')
        } else {
          $n
            .parents('.photos.singleton')
            .toggleClass('singleton gallery')
          $n.remove()
        }
      } else if (!$n
        .find('input.photo')
        .attr('disabled', !$n.hasClass('editing'))
        .attr('disabled')
      ) {
        $n.find('input.photo').focus()
      }
    })
    .on('click', '>.rows>.row>div.commit', e => {
      e.stopPropagation()

      let $selected = $(e.currentTarget).parents('.row').first()
      let url = `/photos/${$(e.delegateTarget).data('owner')}`
      let other

      if ($selected.hasClass('adding')) {
        let file = $selected.find('>.photo').get(0).files[0]

        other = {
          method: 'POST',
          data: ((result = new FormData()) => {
            result.append('file', file, file.name)
            return result
          })(),
          error: _ => { $selected.remove() },
          processData: false,
          contentType: false,
        }
      } else if ($selected.hasClass('editing')) {
        let file = $selected.find('>.photo').get(0).files[0]

        other = {
          url: `${url}/${$selected.attr('id')}`,
          method: 'PATCH',
          data: ((result = new FormData()) => {
            result.set('file', file, file.name)
            return result
          })(),
          processData: false,
          contentType: false,
        }
      } else {
        other = {
          url: `${url}/${$selected.attr('id')}`,
          method: 'DELETE',
          success: (result, status, xhr) => {
            if ($selected.next().click().length === 0) {
              $selected.prev().click()
            }

            $(e.delegateTarget).toggleClass('gallery singleton')

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
              .find('input')
              .attr('disabled', true)
          }
        },
        ...other
      })
    })
    .on('click', '>.rows>.row>img', e => {
      if ($(e.delegateTarget).hasClass('gallery')) {
        $(e.currentTarget)
          .parents('.row')
          .first()
          .toggleClass('selected')
          .find('.toggle-notes')
          .click()
          .parents('.photos')
          .first()
          .toggleClass('gallery singleton')
      } else if ($(e.delegateTarget).hasClass('singleton')) {
        $('body>div.imageview').trigger('activate', $(e.currentTarget).attr('src'))
      }
    })
    .on('click', '>div.add-photo', e => {
      e.stopPropagation()

      let $row = $rowtmpl
        .clone(true, true)
        .toggleClass('template removable editing adding selected')
        // .insertBefore($(e.delegateTarget)
        //   .find('>.rows')
        //   .children()
        //   .first())
        .prependTo($(e.delegateTarget).find('>.rows'))
        .trigger('send')

      // when i find out who's removing active from the template there'll be hell to pay
      $row.find('>.photo-ctl').addClass('active')

      $row
        .find('input.photo')
        .attr('disabled', false)
        .click()
        .parents('.photos')
        .first()
        .toggleClass('gallery singleton')
        .find('.notes>.rows>.row>div.toggle-notes')
        .click()
    })

  // $(window).on('paste', e => {
  //   $('.rows>.row>div.paste-file')
  //     .parents('.row')
  //     .first()
  //     .find('input.photo')
  //     .get(0)
  //     .files = e.clipboardData.files;
  // })
})
