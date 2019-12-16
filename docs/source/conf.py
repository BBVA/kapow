# Configuration file for the Sphinx documentation builder.
#
# This file only contains a selection of the most common options. For a full
# list see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Path setup --------------------------------------------------------------

# If extensions (or modules to document with autodoc) are in another directory,
# add these directories to sys.path here. If the directory is relative to the
# documentation root, use os.path.abspath to make it absolute, like shown here.
#
# import os
# import sys
# sys.path.insert(0, os.path.abspath('.'))

# -- Project information -----------------------------------------------------

project = 'Kapow!'
copyright = '2019, BBVA Innovation Labs'
author = 'BBVA Innovation Labs'


# -- General configuration ---------------------------------------------------

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
extensions = [
    'sphinx.ext.todo',
    'sphinx.ext.imgconverter'
]

# Add any paths that contain templates here, relative to this directory.
templates_path = ['_templates']

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
# This pattern also affects html_static_path and html_extra_path.
exclude_patterns = []

rst_prolog = """
.. role:: tech(code)
   :class: xref

.. role:: nref-option(code)
   :class: xref

.. default-role:: tech

"""

# -- Options for HTML output -------------------------------------------------

# The theme to use for HTML and HTML Help pages.  See the documentation for
# a list of builtin themes.
#
html_theme = "sphinx_rtd_theme"
html_logo = "_static/logo-200px.png"
html_theme_options = {
    'logo_only': True,
    'collapse_navigation': False,
    'navigation_depth': 3,
    'includehidden': True,
    'titles_only': False

}

# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ['_static']

# https://stackoverflow.com/a/56448499
master_doc = 'index'

latex_logo = '_static/logo.png'
latex_documents = [
    ('latextoc',
     'kapow.tex',
     'Kapow! Documentation',
     'BBVA Innovation Labs',
     'manual',
     True)
]

man_pages = [
    ('concepts/resource_tree',
     'kapow-resources',
     'Kapow! Resource Tree Reference',
     'BBVA Innovation Labs',
     1),
    ('examples/examples',
     'kapow-examples',
     'Kapow! Usage Examples',
     'BBVA Innovation Labs',
     1),
]
