$(function () {
  var $buttonbars = $('body>.main>.workspace .buttonbar')
  $buttonbars
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
          .find(`>.${clazz}`)
          .on('click', wrapper)
          .addClass('active')

        wrappers[clazz] = wrapper
      }
    })
    .append($('<img class="remove active" />'))
    .append($('<img class= "ok" />'))
    .append($('<img class="cancel" />'))
    .append($('<img class="add active" />'))
    .append($('<img class="edit active" />'))
    .append($('<img class="refresh active" />'))

  $buttonbars.find('img').attr('src', '/images/transparent.png')
})