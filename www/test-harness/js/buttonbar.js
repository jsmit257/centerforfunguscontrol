$(() => {
  $('.buttonbar')
    .on('reset', (e, h) => {
      var $buttonbar = $(e.currentTarget)
      $buttonbar
        .find('>.add, >.edit')
        .addClass('active')

      $buttonbar
        .find('>.ok, >.cancel')
        .removeClass('active')

      for (let [clazz, handler] of Object.entries(h || {})) {
        $buttonbar.find(`.${clazz}`).off('click', handler)
      }
    })
    .on('set', (e, data) => {
      var $buttonbar = $(e.currentTarget)
      var wrappers = {}

      $buttonbar
        .find('>.add, >.edit, >.ok, >.cancel')
        .removeClass('active')

      for (let [clazz, handler] of Object.entries(data.handlers)) {
        let wrapper = (e) => {
          handler(e)
          data.target.find('.row.selected').removeClass('editing')
          $buttonbar.trigger('reset', wrappers)
        }

        $buttonbar
          .find(`>.${clazz}, >.optional>.${clazz}`)
          .on('click', wrapper)
          .addClass('active')

        wrappers[clazz] = wrapper
      }
    })
    .on('subscribe', (e, cfg) => {
      $(new Image())
        .addClass(`${cfg.clazz} active button`)
        .attr('src', '/images/transparent.png')
        .appendTo($(e.currentTarget).find('>div.optional'))
        .on('click', e => {
          if (!$(e.currentTarget).hasClass('active')) {
            return
          }
          cfg.clicker(e)
        })
    })
    .append($('<div class="optional">')
      .append($('<img class="remove button" />'))) // .html('&#xe020;')
    .append($('<img class="cancel button" />'))
    .append($('<img class= "ok button" />')) // .html('&#xe013;')
    .append($('<img class="add active button" />'))
    .append($('<img class="edit button" />'))
    .append($('<img class="refresh active button" />'))
    .find('img')
    .attr('src', '/images/transparent.png')
})